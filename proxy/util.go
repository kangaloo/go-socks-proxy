package proxy

import (
	"github.com/kangaloo/go-socks-proxy/monitor"
	"log"
	"net"
)

// forward函数退出前需要关闭两个conn，并删除 counter
func forward(dst, src net.Conn, counter *monitor.Counter) {

	defer func() {
		// todo 只关闭一个，以解决下面的 read closed connection 问题，需要考虑关闭哪个
		_ = dst.Close()
		_ = src.Close()
		//counter.Close()
		counter.Done()
	}()

	buf := make([]byte, bufSize)

	for {

		// 此处的错误目前发现两种 io.EOF | read closed connection
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
