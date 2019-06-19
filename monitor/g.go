package monitor

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var registerLock = &sync.Mutex{}

// todo 修改监控项注册策略，每个源地址:目的地址对，注册一个监控项
//  使用独立的register，不实用DefaultRegister

<<<<<<< HEAD
var globalRegisterLock = &sync.Mutex{}

=======
>>>>>>> dev
func Prometheus() {
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	err := http.ListenAndServe("[::]:8080", http.DefaultServeMux)
	log.Panic(err)
}

// todo 增加引用计数和超时字段
//  引用计数为0，且等待了一段时间后，关闭 注销
//  引用计数用 sync.WaitGroup 实现
type Counter struct {
	sync.Mutex
	sync.WaitGroup
	Src        string
	Dst        string
	FlowType   string
	TotalBytes int
	Ch         chan int
	Desc       *prometheus.Desc
	finished   bool
	fresh      bool
}

func NewFlowCounter(src, dst, flow string) *Counter {
	collector := &Counter{
		TotalBytes: 0,
		Ch:         make(chan int, 1024),
		Src:        src,
		Dst:        dst,
		FlowType:   flow,
		Desc: prometheus.NewDesc(
			"SOCKS_PROXY",
			"GO-SOCKS-PROXY TOTAL FLOW",
			nil,
			map[string]string{"SRC": src, "DST": dst, "FlowType": flow},
		),
	}

	registerLock.Lock()
	err := prometheus.Register(collector)
	// 增加引用计数也在这个位置
	registerLock.Unlock()

	if err == nil {
		// todo 如果metric已经注册过，则不启动counter
		go collector.Count()
		collector.Add(1)
		return collector
	}

	// todo 对以下方法进行测试
	exCollector, ok := err.(prometheus.AlreadyRegisteredError)
	if ok {
		log.Info("the error is AlreadyRegisteredError")
		exc, ok := exCollector.ExistingCollector.(*Counter)
		if ok {
			log.Info("convert successful")
			fmt.Printf("%#v", exc)
			exc.Add(1)
			return exc
		}
	}

	// todo 如果metric已经注册过，则不启动counter
	go collector.Count()
	collector.Add(1)
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

	log.Printf("counter for %s exited", c.Dst)
}

func (c *Counter) UnRegister() {
	registerLock.Lock()

	// todo 在这个位置检查引用计数
	// todo 在这个位置 关闭通道
	log.Printf("UnRegister metric successfully: %v", prometheus.Unregister(c))
	registerLock.Unlock()
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

	registerLock.Lock()
	if c.finished {
		c.UnRegister()
	}

	registerLock.Unlock()
	return count
}

func (c *Counter) Release() {
	c.Done()
}

func (c *Counter) Close() {
	//panic: close of closed channel
	// 需要修改，不能马上close 不能马上注销 其他的还要用
	close(c.Ch)
}

func (c *Counter) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Desc
}

func (c *Counter) Collect(ch chan<- prometheus.Metric) {
	count := c.Read()
	counter, err := prometheus.NewConstMetric(
		c.Desc,
		prometheus.CounterValue,
		float64(count),
	)

	if err != nil {
		log.Warn(err)
	}
	ch <- counter
}

type CollectionManager struct{}
