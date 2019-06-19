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
		//log.Printf("%d bytes read, %s, %v", n, buf[:n], buf[:n])
		log.Printf("read %d bytes from %s\n", n, src.RemoteAddr().String())

		if err != nil {
			log.Printf("%#v", err)
			log.Printf("read forward: %s", err)
			return
		}

		counter.Write(n)

		_, err = dst.Write(buf[:n])
		if err != nil {
			log.Printf("%#v", err)
			log.Printf("write forward: %s", err)
			return
		}
	}
}
