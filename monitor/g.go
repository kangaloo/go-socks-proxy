package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

// todo 修改监控项注册策略，每个源地址:目的地址对，注册一个监控项
//  使用独立的register，不实用DefaultRegister

var globalRegisterLock = &sync.Mutex{}

func Prometheus() {
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	err := http.ListenAndServe("[::]:8080", http.DefaultServeMux)
	log.Panic(err)
}

type Counter struct {
	sync.Mutex // for compute and read total concurrent
	prometheus.Desc
	Src         string
	Dst         string
	FlowType    string
	TotalBytes  int
	Ch          chan int
	finished    bool
	fresh       bool
	invokeCount int
}

func (c *Counter) Add() {

	c.invokeCount += 1
}

func (c *Counter) Done() {
	c.invokeCount -= 1
}

func NewFlowCounter(src, dst, flow string) *Counter {
	collector := &Counter{
		Desc: *prometheus.NewDesc(
			"SOCKS_PROXY_FLOW",
			"SOCKS-PROXY TOTAL FLOW",
			nil,
			map[string]string{"SRC": src, "DST": dst, "FlowType": flow},
		),
		Ch:       make(chan int, 1024),
		Src:      src,
		Dst:      dst,
		FlowType: flow,
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
				log.Info("convert ExistingCollector to `*Counter` successful")
				existErrCollector.Add()
				globalRegisterLock.Unlock()
				return existErrCollector
			}
		}
		log.Warn("convert err to prometheus.AlreadyRegisteredError failed")
		log.Panic(err)
	}

	globalRegisterLock.Unlock()
	collector.Add()
	go collector.Count()
	return collector
}

func (c *Counter) Count() {
	for {
		n, ok := <-c.Ch
		if !ok {
			break
		}

		log.Printf("get %d from channel\n", n)
		c.Lock()
		c.fresh = false
		c.TotalBytes += n
		c.Unlock()
	}
	c.finished = true
	log.Info("counter for %s exited", c.Dst)
}

func (c *Counter) UnRegister() {
	globalRegisterLock.Lock()
	if c.invokeCount >= 0 {
		globalRegisterLock.Unlock()
		return
	}
	close(c.Ch)
	log.Printf("UnRegister metric successfully: %v", prometheus.Unregister(c))
	globalRegisterLock.Unlock()
}

func (c *Counter) Write(n int) {
	c.Ch <- n
	log.Debug("net process write %d to collector channel\n", n)
}

func (c *Counter) Read() int {
	c.Lock()
	count := c.TotalBytes
	c.TotalBytes = 0
	c.fresh = true
	c.Unlock()
	log.Debug("prometheus read %d from collector", count)

	/*
		registerLock.Lock()
		if c.finished {
			c.UnRegister()
		}

		registerLock.Unlock()

	*/
	return count
}

func (c *Counter) Describe(ch chan<- *prometheus.Desc) {
	ch <- &c.Desc
}

func (c *Counter) Collect(ch chan<- prometheus.Metric) {
	count := c.Read()
	counter, err := prometheus.NewConstMetric(
		&c.Desc,
		prometheus.CounterValue,
		float64(count),
	)

	if err != nil {
		log.Warn(err)
	}
	ch <- counter
	c.UnRegister()
}
