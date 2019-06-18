package proxy

import (
	"github.com/kangaloo/go-socks-proxy/monitor"
	"log"
	"net"
)

// forward函数退出前需要关闭两个conn，并删除 counter
// todo 怀疑出错的地方在这里，两个协程互相影响
func forward(dst, src net.Conn, counter *monitor.Counter) {

	defer func() {
		_ = dst.Close()
		_ = src.Close()
		counter.Close()
	}()

	buf := make([]byte, bufSize)

	for {
		n, err := src.Read(buf)
		log.Printf("%d bytes read, %s, %v", n, buf[:n], buf[:n])

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
