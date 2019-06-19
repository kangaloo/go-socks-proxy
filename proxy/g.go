package proxy

import (
	"github.com/kangaloo/go-socks-proxy/protocol"
	log "github.com/sirupsen/logrus"
	"net"
)

var bufSize = 1024

func SimpleForward(conn net.Conn) {
	// 此处的连接不能关闭 向客户端返回数据需要这个连接
	/*
		defer func() {
			log.Printf("connection from %s will be closed\n", conn.RemoteAddr())
			_ = conn.Close()
		}()

	*/

	conn, addr, err := protocol.Socks(conn)
	if err != nil {
		log.Printf("%#v", err)
		_ = conn.Close()
		return
	}

	TCPProxy(conn, addr)
}
