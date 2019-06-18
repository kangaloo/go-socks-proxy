package proxy

import (
	"github.com/kangaloo/go-socks-proxy/monitor"
	"log"
	"net"
)

// forward函数退出前需要关闭两个conn，并删除 counter
func forward(dst, src net.Conn, counter *monitor.Counter) {

	defer func() {
		_ = dst.Close()
		_ = src.Close()
		counter.Close()
	}()

	buf := make([]byte, bufSize)

	for {
		n, err := src.Read(buf)
		if err != nil {
			log.Printf("%#v", err)
			return
		}

		_, err = dst.Write(buf[:n])
		if err != nil {
			log.Printf("%#v", err)
			return
		}
	}
}
