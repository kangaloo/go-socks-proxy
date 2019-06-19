package main

import (
	_ "github.com/kangaloo/go-socks-proxy/logger"
	"github.com/kangaloo/go-socks-proxy/monitor"
	"github.com/kangaloo/go-socks-proxy/proxy"
	log "github.com/sirupsen/logrus"
	"net"
)

// 加入 prometheus 监控
// socks 代理

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

		log.Printf("new connection from %s\n", conn.RemoteAddr())
		go proxy.SimpleForward(conn)
		//go communicate(conn)
	}
}

// 每个communicate协程传入一个channel，用于写入转发的byte数，prom监控读取channel中的数据，实现统计流量

/*
func communicate(conn net.Conn) {
	bufSize := 8
	buf := make([]byte, bufSize)

	defer func() {
		log.Printf("connection from %s will be closed\n", conn.RemoteAddr())
		_ = conn.Close()
	}()

	_, _ = conn.Write([]byte("socks-proxy >> "))

	counter := monitor.NewFlowCounter(conn.RemoteAddr().String())
	defer counter.Close()

	for {
		// client 主动关闭时且无数据可读时读到了EOF (TCP收到了断开连接的包)
		// 当还有没读完的数据，client主动关闭时会怎么样 (猜测会先处理完正常数据包，在处理断开连接的包，因为是流式的)
		// net包对底层做了封装，在TCP连接处于某种状态时向上层返回EOF，比如TCP连接收到断开指令时
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("%#v", err)
			log.Println(err)
			return
		}

		//monitor.FlowCounter.Write(n)
		counter.Write(n)

		fmt.Printf("%s, %d\n", buf[:n], n)
		if n < bufSize {
			_, _ = conn.Write([]byte("socks-proxy >> "))
			//break
		}
	}
}
*/
