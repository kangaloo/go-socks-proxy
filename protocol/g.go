package protocol

const (
	methodWithoutAuth byte = iota
	methodAuth
)

const methodNotSupport byte = 0xff

const (
	addrTypeIPv4       byte = 0x01
	addrTypeDomainName byte = 0x03
	addrTypeIPv6       byte = 0x04
)
