package proxy

import (
	"github.com/kangaloo/go-socks-proxy/monitor"
	"log"
	"net"
	"time"
)

// 在该函数中完成prometheus监控指标的注册
func TCPProxy(srcConn net.Conn, dstAddr string) {
	// timeout 20s
	dstConn, err := net.DialTimeout("tcp", dstAddr, time.Second*20)
	if err != nil {
		log.Printf("%#v", err)
		return
	}

	// 监控指标需要更改，需要 in/out 标签，每个指标有两个连接，是否都监听
	// 监控指标 只需要 源地址标签 和 目的地址标签 不需要端口 没有必要
	// 为解决 metric 的重复注册问题还是需要端口号，在监控端做数据聚合，让端口不同但源、目的一样的流量统计在一起
	flowInCounter := monitor.NewFlowCounter(srcConn.RemoteAddr().String())
	go forward(dstConn, srcConn, flowInCounter)

	flowOutCounter := monitor.NewFlowCounter(dstConn.RemoteAddr().String())
	go forward(srcConn, dstConn, flowOutCounter)
}
