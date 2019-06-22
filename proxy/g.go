package proxy

import (
	"github.com/kangaloo/go-socks-proxy/protocol"
	log "github.com/sirupsen/logrus"
	"net"
)

var bufSize = 1024

func SimpleForward(conn net.Conn) {
	conn, addr, err := protocol.Socks(conn)
	if err != nil {
		//log.Printf("%#v", err)
		log.Warn(err)
		_ = conn.Close()
		return
	}

	TCPProxy(conn, addr)
}
