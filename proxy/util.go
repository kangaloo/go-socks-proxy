package proxy

import (
	"github.com/kangaloo/go-socks-proxy/monitor"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
	"sync"
)

// forward函数退出前需要关闭两个conn，并删除 counter
func forward(id string, gw *sync.WaitGroup, dst, src net.Conn, counter *monitor.Counter) {

	fields := log.Fields{
		"proxies_id": id,
		"conn":       src.RemoteAddr().String() + "->" + dst.RemoteAddr().String(),
	}

	defer func() {
		// todo 只关闭一个，以解决下面的 read closed connection 问题，需要考虑关闭哪个
		//_ = dst.Close()
		log.WithFields(fields).Info(src.Close())
		counter.Done()
		gw.Done()
	}()

	buf := make([]byte, bufSize)

	for {
		// 此处的错误目前发现两种 io.EOF | read closed connection
		// 修改后应该不会在出现read closed connection
		n, err := src.Read(buf)
		//log.Printf("%d bytes read, %s, %v", n, buf[:n], buf[:n])

		// 新出现的错误 `read tcp 10.0.0.73:1080->111.164.175.28:53368: read: connection reset by peer`
		if err != nil {
			log.WithFields(fields).Warn(err)
			return
		}

		log.WithFields(fields).Printf("read %d bytes from %s", n, src.RemoteAddr().String())

		counter.Write(n)
		_, err = dst.Write(buf[:n])
		if err != nil {
			log.WithFields(fields).Warn(err)
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
