package proxy

import (
	"github.com/kangaloo/go-socks-proxy/protocol"
	"log"
	"net"
)

var bufSize = 1024

func SimpleForward(conn net.Conn) {

	// todo 此处的连接不能关闭，需要修改
	defer func() {
		log.Printf("connection from %s will be closed\n", conn.RemoteAddr())
		_ = conn.Close()
	}()

	conn, addr, err := protocol.Socks(conn)
	if err != nil {
		log.Printf("%#v", err)
		return
	}

	TCPProxy(conn, addr)

}
