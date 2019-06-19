package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"sync"
)

// todo 修改监控项注册策略，每个源地址:目的地址对，注册一个监控项
//  使用独立的register，不实用DefaultRegister

func Prometheus() {
	/*
		counter := &Counter{
			Zone: "zone",
			Flow: prometheus.NewDesc(
				"socks_proxy",
				"GO-SOCKS-PROXY TOTAL FLOW",
				[]string{"proxy"},
				map[string]string{"zone": "zone"},
			),
		}

		err := prometheus.Register(counter)
		if err != nil {
			log.Println(err)
		}

	*/

	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	err := http.ListenAndServe("[::]:8080", http.DefaultServeMux)
	log.Panic(err)
}

type Counter struct {
	sync.Mutex
	Zone       string
	Address    string
	TotalBytes int
	Ch         chan int
	Flow       *prometheus.Desc
	finished   bool
}

func NewFlowCounter(addr string) *Counter {
	counter := &Counter{
		Zone:       "zone",
		Address:    addr,
		TotalBytes: 0,
		Ch:         make(chan int, 1024),
		Flow: prometheus.NewDesc(
			"SOCKS_PROXY",
			"GO-SOCKS-PROXY TOTAL FLOW",
			[]string{"flow", "proxy"},
			map[string]string{"ADDR": addr},
		),
	}
	err := prometheus.Register(counter)
	if err != nil {
		log.Println(err)
	}

	go counter.Count()
	return counter
}

func (c *Counter) Count() {
	for {
		n, ok := <-c.Ch
		if !ok {
			break
		}

		log.Printf("get %d from channel\n", n)
		c.Lock()
		c.TotalBytes += n
		c.Unlock()
	}
	c.finished = true

	log.Printf("counter for %s exited", c.Address)
}

func (c *Counter) UnRegister() {
	log.Printf("UnRegister metric successfully: %v", prometheus.Unregister(c))
}

func (c *Counter) Write(n int) {
	c.Ch <- n
	log.Printf("write %d bytes, the total is %d bytes\n", n, c.TotalBytes)
}

func (c *Counter) Read() int {
	c.Lock()
	count := c.TotalBytes
	//c.TotalBytes = 0
	c.Unlock()
	log.Printf("prometheus read %d bytes", count)
	if c.finished {
		c.UnRegister()
	}
	return count
}

func (c *Counter) Close() {
	close(c.Ch)
}

func (c *Counter) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Flow
}

func (c *Counter) Collect(ch chan<- prometheus.Metric) {
	/*
		count := FlowCounter.Read()
		counter, err := prometheus.NewConstMetric(
			c.Flow,
			prometheus.CounterValue,
			float64(count),
			"flow",
		)

		if err != nil {
			log.Println(err)
		}

		ch <- counter

	*/

	count := c.Read()
	counter, err := prometheus.NewConstMetric(
		c.Flow,
		prometheus.CounterValue,
		float64(count),
		"",
		"",
	)

	if err != nil {
		log.Println(err)
	}

	ch <- counter
}
