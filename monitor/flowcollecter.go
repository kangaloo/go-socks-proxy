package monitor

import (
	"log"
	"sync"
)

// 改用channel，实现无锁操作
type FlowCollector struct {
	sync.Mutex
	Ch chan int
	TotalBytes int
}

func (c *FlowCollector) Write(n int)  {
	c.Ch <- n
	log.Printf("write %d bytes, the total is %d bytes\n", n, c.TotalBytes)
}

func (c *FlowCollector) Read() int {
	c.Lock()
	count := c.TotalBytes
	c.TotalBytes = 0
	c.Unlock()
	log.Printf("read %d bytes", count)
	return count
}

func (c *FlowCollector) Count() {
	for {
		n := <- c.Ch
		log.Printf("get %d from channel\n", n)
		c.Lock()
		c.TotalBytes += n
		c.Unlock()
	}
}

var FlowCounter = &FlowCollector{Ch: make(chan int, 1024)}

func init()  {
	go FlowCounter.Count()
}