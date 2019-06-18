package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
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

func Socks(conn net.Conn) (net.Conn, string, error) {
	//defer func() {log.Printf("close connection %s: %#v", conn.RemoteAddr(), conn.Close())}()
	bufSize := 256
	buf := make([]byte, bufSize)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("%#v", err)
		return conn, "", err
	}

	log.Printf("read protocol handshake packet %v", buf[:n])
	if !sliceEqual(socksMethodPacket, buf[:n]) {
		return conn, "", errors.New("packet from client Unrecognized")
	}

	_, err = conn.Write([]byte{0x05, 0x00})
	if err != nil {
		log.Printf("%#v", err)
		return conn, "", err
	}

	n, err = conn.Read(buf)
	if err != nil {
		log.Printf("%#v", err)
		return conn, "", err
	}
	log.Printf("read protocol handshake packet %v", buf[:n])

	addr, err := parseAddr(buf[:n])
	log.Printf("%s, %s", conn.RemoteAddr(), addr)

	_, _ = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	return conn, addr, err
}

// todo 该函数仅支持域名，需要补充IP地址的parser和UDP的parser
func parseAddr(packet []byte) (string, error) {
	// IPv4 4 bytes, port 2 bytes, socks 4 bytes
	if len(packet) < 9 {
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

	log.Printf("parse port from packet: %d", port)

	// todo 这个位置为什么是5
	log.Printf("addr slice: %#v", packet[5:len(packet)-2])

	addr := string(packet[5 : len(packet)-2])
	log.Printf("get addr: %#v\n", addr)

	return addr + ":" + strconv.Itoa(int(port)), nil
}
