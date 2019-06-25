package monitor

import (
	"github.com/kangaloo/go-socks-proxy/util"
)

// IDGenerator is a global id generator
var IDGenerator *util.Generator

func init() {
	IDGenerator = &util.Generator{}
}
