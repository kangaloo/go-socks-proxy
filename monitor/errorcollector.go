package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	PacketErr = &errorCollector{
		fqName:     "SOCKS_PROXY_FAILED_TOTAL",
		help:       "SOCKS_PROXY_FAILED_TOTAL",
		constLabel: map[string]string{"ErrorType": "Packet Unrecognized"},
	}
	DialErr = &errorCollector{
		fqName:     "SOCKS_PROXY_FAILED_TOTAL",
		help:       "SOCKS_PROXY_FAILED_TOTAL",
		constLabel: map[string]string{"ErrorType": "Dial To Destination Address Failed"},
	}
	SocksErr = &errorCollector{
		fqName:     "SOCKS_PROXY_FAILED_TOTAL",
		help:       "SOCKS_PROXY_FAILED_TOTAL",
		constLabel: map[string]string{"ErrorType": "Socks Protocol Error"},
	}

	// 记录monitor包出现的问题
	monitorErr = &errorCollector{
		fqName:     "SOCKS_PROXY_FAILED_TOTAL",
		help:       "SOCKS_PROXY_FAILED_TOTAL",
		constLabel: map[string]string{"ErrorType": "Monitor Error"},
	}
	// 记录monitor包出现的warning
	monitorWarn = &errorCollector{
		fqName:     "SOCKS_PROXY_FAILED_TOTAL",
		help:       "SOCKS_PROXY_FAILED_TOTAL",
		constLabel: map[string]string{"ErrorType": "Monitor Warning"},
	}
)

func init() {
	PacketErr.Setup()
	DialErr.Setup()
	SocksErr.Setup()
	monitorErr.Setup()
	monitorWarn.Setup()
}

// 外部使用时 monitor.PacketErr.Write(1, src, dst)
// src dst 两个字段需要再考虑 暂时不增加这两个字段

// errorCollector is used by count errors occurring during forwarding
type errorCollector struct {
	sync.Mutex // 并发写入是使用
	prometheus.Desc
	fqName     string
	help       string
	constLabel map[string]string
	count      int
}

func (c *errorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- &c.Desc
}

func (c *errorCollector) Collect(ch chan<- prometheus.Metric) {
	count := c.read()
	metric, err := prometheus.NewConstMetric(
		&c.Desc,
		prometheus.CounterValue,
		float64(count),
	)
	if err != nil {
		log.Warn(err)
	}
	ch <- metric
}

func (c *errorCollector) Setup() {
	c.Desc = *prometheus.NewDesc(
		c.fqName,
		c.help,
		nil,
		c.constLabel,
	)
	err := prometheus.Register(c)
	if err != nil {
		log.Warn(err)
	}
}

//func (c *errorCollector) SetLabel(label string) {}

// 提供给记录错误的地方使用
func (c *errorCollector) Write(n int) {
	c.Lock()
	defer c.Unlock()
	c.count += n
}

func (c *errorCollector) read() int {
	c.Lock()
	defer c.Unlock()
	return c.count
}
