package proxy

import (
	"github.com/kangaloo/go-socks-proxy/monitor"
	"github.com/kangaloo/go-socks-proxy/protocol"
	log "github.com/sirupsen/logrus"
	"net"
)

var bufSize = 1024

func SimpleForward(conn net.Conn) {
	conn, dst, err := protocol.Socks(conn)
	if err != nil {
		// 记录一条socks protocol包出现的错误
		monitor.SocksErr.Write(1)
		log.Warn(err)
		_ = conn.Close()
		return
	}

	proxies := NewProxies(conn, dst)
	log.WithField("request_domain", dst.RemoteAddr()).Infof("create new proxy connection successfully")
	proxies.Run()
}
