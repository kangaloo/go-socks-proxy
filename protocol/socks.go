package protocol

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

// proxy函数需要访问destAddr，并将返回的数据写入Socks返回的conn，
// 中间过程中加入流量监控

type packetError struct {
	message string
	packet  []byte
}

func (e *packetError) Error() string {
	return fmt.Sprintf("%s, packet: %v", e.message, e.packet)
}

// Socks
func Socks(conn net.Conn) (net.Conn, net.Conn, error) {
	// 读取客户端发送的第一个数据包
	packet, err := socksPacket(conn)
	if err != nil {
		return conn, nil, err
	}

	log.Debug("read protocol handshake packet %v", packet)

	method, err := consultMethod(packet)
	if err != nil {
		// send method not supported to client
		_, _ = conn.Write(genPacket(method))
		return conn, nil, err
	}

	_, err = conn.Write(genPacket(method))
	if err != nil {
		return conn, nil, err
	}

	// 提取目的地址
	dst, err := extractDstAddr(conn)
	if err != nil {
		// 0x08: Address type not supported
		_, _ = conn.Write([]byte{0x05, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return conn, nil, err
	}

	dstConn, err := net.DialTimeout("tcp", dst, time.Second*10)
	if err != nil {
		_, _ = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return conn, nil, err
	}

	// todo 向客户端返回代理服务器使用的ip和port
	_, _ = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	return conn, dstConn, err
}
