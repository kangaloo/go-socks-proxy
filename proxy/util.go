package proxy

import (
	"github.com/kangaloo/go-socks-proxy/monitor"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
)

// forward函数退出前需要关闭两个conn，并删除 counter
func forward(dst, src net.Conn, counter *monitor.Counter) {

	defer func() {
		// todo 只关闭一个，以解决下面的 read closed connection 问题，需要考虑关闭哪个
		_ = dst.Close()
		//_ = src.Close()
		counter.Done()
	}()

	buf := make([]byte, bufSize)

	for {

		// 此处的错误目前发现两种 io.EOF | read closed connection
		n, err := src.Read(buf)
		//log.Printf("%d bytes read, %s, %v", n, buf[:n], buf[:n])
		log.Printf("read %d bytes from %s", n, src.RemoteAddr().String())

		if err != nil {
			//log.Printf("%#v", err)
			log.Warn(err)
			return
		}

		counter.Write(n)

		_, err = dst.Write(buf[:n])
		if err != nil {
			//log.Printf("%#v", err)
			log.Warn(err)
			return
		}
	}
}

func addrFormat(addr string) string {
	addr = strings.Split(addr, ":")[0]
	// todo 需要检查，如果是合法的IP，则不做处理
	//addr = strings.Join(strings.Split(addr, ".")[1:], ".")
	// todo 检查截断后的域名是否合法，如 github.com 截断后为 com 是不正确的
	return addr
}
