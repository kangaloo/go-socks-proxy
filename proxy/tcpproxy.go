package proxy

import (
	"github.com/kangaloo/go-socks-proxy/monitor"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)

// 在该函数中完成prometheus监控指标的注册
func TCPProxy(srcConn net.Conn, dstAddr string) {
	// timeout 20s, 需要变成可配置的项
	dstConn, err := net.DialTimeout("tcp", dstAddr, time.Second*10)
	if err != nil {
		monitor.DialErr.Write(1)
		log.WithFields(log.Fields{
			"request_domain": dstAddr,
			"client_address": srcConn.RemoteAddr().String(),
		}).Warn(err)
		return
	}

	log.WithField("request_domain", dstAddr).Infof("create new proxy connection successfully")

	// 监控指标需要更改，需要 in/out 标签，每个指标有两个连接，是否都监听
	// 监控指标 只需要 源地址标签 和 目的地址标签 不需要端口 没有必要
	// 为解决 metric 的重复注册问题还是需要端口号，在监控端做数据聚合，让端口不同但源、目的一样的流量统计在一起

	// SOCKS_PROXY{src=clientIP, dst=baidu.com, flow_type=inflow}

	// todo 新增一个处理函数，处理域名，只取域名后缀作为DST

	// in
	flowOutCounter := monitor.NewFlowCounter(strings.Split(srcConn.RemoteAddr().String(), ":")[0], addrFormat(dstAddr), "download")
	go forward(srcConn, dstConn, flowOutCounter)

	// SOCKS_PROXY{src=clientIP, dst=baidu.com, flow_type=outflow}
	// out
	flowInCounter := monitor.NewFlowCounter(strings.Split(srcConn.RemoteAddr().String(), ":")[0], addrFormat(dstAddr), "upload")
	go forward(dstConn, srcConn, flowInCounter)
}
