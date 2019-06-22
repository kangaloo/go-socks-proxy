package main

import (
	_ "github.com/kangaloo/go-socks-proxy/logger"
	"github.com/kangaloo/go-socks-proxy/monitor"
	"github.com/kangaloo/go-socks-proxy/proxy"
	log "github.com/sirupsen/logrus"
	"net"
)

func main() {
	go monitor.Prometheus()

	server, err := net.Listen("tcp", "[::]:1080")
	if err != nil {
		log.Panic(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println(err)
		}

		log.Printf("new connection from %s", conn.RemoteAddr())
		go proxy.SimpleForward(conn)
	}
}
