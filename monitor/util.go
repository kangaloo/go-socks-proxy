package monitor

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"sync"
	"time"
)

type Generator struct {
	sync.Mutex
	seed int
}

func (g *Generator) Generate() string {
	g.Lock()
	defer g.Unlock()
	g.seed += 1
	h := md5.New()
	h.Write([]byte(strconv.Itoa(int(time.Now().UnixNano()))))
	id := hex.EncodeToString(h.Sum([]byte(strconv.Itoa(g.seed))))
	return id
}

// IDGenerator is a global id generator
var IDGenerator *Generator

func init() {
	IDGenerator = &Generator{}
}
