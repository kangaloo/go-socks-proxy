package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

// 生成返回给客户端的包
func genPacket(method byte) []byte {
	return []byte{0x05, method}
}

// 协商method
func consultMethod(packet []byte) (byte, error) {
	if err := checkVersion(packet); err != nil {
		return methodNotSupport, err
	}
	if err := checkLen(packet); err != nil {
		return methodNotSupport, err
	}

	supportedMethods := []byte{methodWithoutAuth, methodAuth}

	for _, clientMethod := range packet[2:] {
		for _, serverMethod := range supportedMethods {
			if clientMethod == serverMethod {
				return clientMethod, nil
			}
		}
	}

	return methodNotSupport, errors.New("methods from client not supported")
}

func checkLen(packet []byte) error {
	if len(packet)-2 == int(packet[1]) {
		return nil
	}
	return errors.New("check length failed")
}

func checkVersion(packet []byte) error {
	if packet[0] == 0x05 {
		return nil
	}
	return errors.New("check version failed")
}

// todo 测试这个函数
// 每次读取一个完成的数据包 仅用于处理socks协议的握手包
func socksPacket(conn net.Conn) ([]byte, error) {

	// 如果数据包的大小刚好为bufSize，第二次循环是是否会阻塞
	// 所以需要按照协议处理数据包，根据协议头获取数据包的大小
	// 不能以是否能继续从socket读到数据为依据
	// 先读取两个字节 版本和method个数
	// 根据方法个继续读取数据

	headerSize := 2
	header := make([]byte, headerSize)

	n, err := conn.Read(header)
	if err != nil {
		return nil, err
	}

	// 网络不好的情况下可能出现无法完整的读取到前两个字节的情况
	if n < 2 {
		return nil, errors.New("read packet error")
	}

	methodSize := int(header[1])
	methods, err := readBytes(conn, methodSize)
	if err != nil {
		return nil, err
	}

	return append(header, methods...), nil
}

// 从conn读取n字节的数据
func readBytes(conn net.Conn, size int) ([]byte, error) {
	bufSize := size
	buf := make([]byte, bufSize)
	packet := make([]byte, 0)

	for {
		size, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}

		packet = append(packet, buf...)
		if len(packet) < size {
			bufSize = size - len(packet)
			buf = buf[:bufSize]
			continue
		}
		break
	}

	return packet, nil
}

// todo 该函数仅支持域名，需要补充IP地址的parser和UDP的parser
func extractDstAddr(conn net.Conn) (string, error) {
	header, err := readBytes(conn, 4)
	if err != nil {
		return "", err
	}

	if err := checkVersion(header); err != nil {
		return "", err
	}

	addrType := header[3]
	if addrType == addrTypeIPv4 {
		// 0x08: Address type not supported
		return "", errors.New("ipv4 address type not supported")
	}

	if addrType == addrTypeDomainName {
		return extractDomain(conn)
	}

	if addrType == addrTypeIPv6 {
		return "", errors.New("ipv6 address type not supported")
	}

	return "", errors.New("unknown address type not supported")
}

func extractDomain(conn net.Conn) (string, error) {
	// get domain length
	packet, err := readBytes(conn, 1)
	if err != nil {
		return "", err
	}
	l := int(packet[0])

	// get domain
	packet, err = readBytes(conn, l)
	if err != nil {
		return "", err
	}

	domain := string(packet)

	// get port
	packet, err = readBytes(conn, 2)
	if err != nil {
		return "", err
	}

	var port int16
	buf := bytes.NewBuffer(packet)
	err = binary.Read(buf, binary.BigEndian, &port)
	if err != nil {
		return "", err
	}

	return domain + ":" + strconv.Itoa(int(port)), nil
}
