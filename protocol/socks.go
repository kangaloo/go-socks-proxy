package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
)

// SocksProxy 实现socks协议
// 只处理协议，处理完协议返回conn
// 返回conn，处理得到的destAddr，error

// proxy函数需要访问destAddr，并将返回的数据写入Socks返回的conn，
// 中间过程中加入流量监控

var (
	socksMethodPacket = []byte{0x05, 0x01, 0x00}
)

type packetError struct {
	message string
	packet  []byte
}

func (e *packetError) Error() string {
	return fmt.Sprintf("%s, packet: %v", e.message, e.packet)
}

func Socks(conn net.Conn) (net.Conn, string, error) {
	bufSize := 256
	buf := make([]byte, bufSize)
	n, err := conn.Read(buf)
	if err != nil {
		return conn, "", err
	}

	log.Debug("read protocol handshake packet %v", buf[:n])
	if !sliceEqual(socksMethodPacket, buf[:n]) {
		return conn, "", &packetError{
			message: "packet from " + conn.RemoteAddr().String() + " Unrecognized",
			packet:  buf[:n],
		}
	}

	_, err = conn.Write([]byte{0x05, 0x00})
	if err != nil {
		return conn, "", err
	}

	n, err = conn.Read(buf)
	if err != nil {
		return conn, "", err
	}

	log.Debug("read protocol handshake packet %v", buf[:n])
	addr, err := parseAddr(buf[:n])
	log.Printf("srcAddr: %s, dstAddr: %s", conn.RemoteAddr(), addr)

	// todo 应该将创建转发连接的工作放在该函数内的这个位置，创建失败向客户端返回错误，才算完整实现了socks协议
	_, _ = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	return conn, addr, err
}

// todo 该函数仅支持域名，需要补充IP地址的parser和UDP的parser
func parseAddr(packet []byte) (string, error) {
	// IPv4 4 bytes, port 2 bytes, socks 4 bytes
	if len(packet) < 10 {
		return "", errors.New("can not parse packet, too short")
	}

	// 地址解析这段代码转移到 util.go 中的 parseDomain 函数中
	var port int16
	buf := bytes.NewBuffer(packet[len(packet)-2:])
	err := binary.Read(buf, binary.BigEndian, &port)
	if err != nil {
		log.Printf("parse port failed: %#v", err)
		return "", err
	}

	log.WithField("port", port).Info("parse port from packet success")

	// 这个位置为什么是5 VER 1, CMD 1, RSV 1, ATYP 1, 地址的第一位是domain长度
	log.Printf("addr slice: %#v", packet[5:len(packet)-2])

	addr := string(packet[5 : len(packet)-2])
	log.Printf("get addr: %#v\n", addr)

	return addr + ":" + strconv.Itoa(int(port)), nil
}
