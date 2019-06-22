package protocol

import (
	"net"
	"testing"
	"time"
)

var server net.Listener

func init() {
	var err error
	server, err = net.Listen("tcp", "[::]:1081")
	if err != nil {
		panic(err)
	}
}

func TestSocks(t *testing.T) {
	conn, err := server.Accept()
	if err != nil {
		t.Error(err)
	}

	_, addr, err := Socks(conn)
	if err != nil {
		t.Error(err)
	}

	proxyConn, err := net.DialTimeout("tcp", addr, time.Second*15)
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		_ = proxyConn.Close()
	}()
}
