package monitor

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"sync"
)

// todo 修改监控项注册策略，每个源地址:目的地址对，注册一个监控项
//  使用独立的register，不实用DefaultRegister

var globalRegisterLock = &sync.Mutex{}

// 初始化全局collector，提供给注册失败且转换失败的collector使用
func init() {}

func Prometheus() {
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	err := http.ListenAndServe("[::]:8080", http.DefaultServeMux)
	log.Panic(err)
}

// Counter is a prometheus metric collector
type Counter struct {
	sync.Mutex // for compute and read total concurrent
	prometheus.Desc
	ID          string
	Src         string
	Dst         string
	FlowType    string
	TotalBytes  int
	Ch          chan int
	finished    bool
	fresh       bool
	invokeCount counter
}

func (c *Counter) Add() {
	c.invokeCount.add()
}

func (c *Counter) Done() {
	c.invokeCount.done()
}

func NewFlowCounter(src, dst, flow string) *Counter {
	collector := &Counter{
		Desc: *prometheus.NewDesc(
			"SOCKS_PROXY_FLOW_BYTES",
			"SOCKS-PROXY TOTAL FLOW",
			[]string{"connections"},
			map[string]string{"SRC": src, "DST": dst, "FlowType": flow},
		),
		Ch:       make(chan int, 1024),
		Src:      src,
		Dst:      dst,
		FlowType: flow,
		ID:       IDGenerator.Generate(),
	}

	globalRegisterLock.Lock()
	err := prometheus.Register(collector)

	if err != nil {
		log.Info("register collector failed, %s, try to convert err to prometheus.AlreadyRegisteredError",
			err.Error())
		existErr, ok := err.(prometheus.AlreadyRegisteredError)

		if ok {
			log.Info("convert err to prometheus.AlreadyRegisteredError successful, " +
				"will try to convert to collector `*Counter`")
			existErrCollector, ok := existErr.ExistingCollector.(*Counter)

			if ok {
				log.WithFields(log.Fields{
					"exist_collector_id": existErrCollector.ID,
					"new_collector_id":   collector.ID,
				}).Info("convert ExistingCollector to `*monitor.Counter` successfully")
				existErrCollector.Add()
				globalRegisterLock.Unlock()
				log.WithField("collector_id", existErrCollector.ID).Info("new connection attached")
				return existErrCollector
			}
		}
		log.Warn("convert err to prometheus.AlreadyRegisteredError failed")
		log.Panic(err)
	}

	// todo 此处需要修改，转换失败时也会返回新的，应该返回一个全局的处理错误的collector
	globalRegisterLock.Unlock()
	collector.Add()
	go collector.Count()
	log.WithField("collector_id", collector.ID).Info("new collector created")
	return collector
}

func (c *Counter) Count() {
	for {
		n, ok := <-c.Ch
		if !ok {
			break
		}

		log.WithField("collector_id", c.ID).Infof("counter get %d from channel", n)
		c.Lock()
		c.fresh = false
		c.TotalBytes += n
		c.Unlock()
	}
	c.finished = true
	log.Infof("count goroutine for %s exited", c.ID)
}

func (c *Counter) UnRegister() {
	globalRegisterLock.Lock()
	log.Debug(c)

	log.WithField("collector_id", c.ID).Infof("collector invoke count is %s", c.invokeCount.string())
	if c.invokeCount.num() > 0 {
		globalRegisterLock.Unlock()
		return
	}
	close(c.Ch)

	if prometheus.Unregister(c) {
		log.WithField("collector_id", c.ID).Info("UnRegister collector successfully")
	} else {
		log.WithField("collector_id", c.ID).Warn("UnRegister collector failed")
	}
	globalRegisterLock.Unlock()
}

func (c *Counter) Write(n int) {
	c.Ch <- n
	log.WithField("collector_id", c.ID).Debugf("net process write %d to collector channel", n)
}

func (c *Counter) Read() int {
	c.Lock()
	count := c.TotalBytes
	c.TotalBytes = 0
	c.fresh = true
	c.Unlock()
	log.Debugf("prometheus read %d from collector", count)
	return count
}

func (c *Counter) String() string {
	return fmt.Sprintf("id: %s, invoke: %d, desc: %s",
		c.ID, c.invokeCount.num(), c.Desc.String())
}

func (c *Counter) Describe(ch chan<- *prometheus.Desc) {
	ch <- &c.Desc
}

func (c *Counter) Collect(ch chan<- prometheus.Metric) {
	count := c.Read()
	metric, err := prometheus.NewConstMetric(
		&c.Desc,
		prometheus.CounterValue,
		float64(count),
		c.invokeCount.string(),
	)

	if err != nil {
		monitorWarn.Write(1)
		log.Warn(err)
	}
	ch <- metric
	c.UnRegister()
}

type counter struct {
	sync.Mutex
	count int
}

func (c *counter) add() {
	c.Lock()
	defer c.Unlock()
	c.count += 1
}

func (c *counter) done() {
	c.Lock()
	defer c.Unlock()
	c.count -= 1
}

func (c *counter) num() int {
	return c.count
}

func (c *counter) string() string {
	return strconv.Itoa(c.count)
}
