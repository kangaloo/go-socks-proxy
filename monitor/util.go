package monitor

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"sync"
	"time"
)

// id generator

// todo 每次产生ID是更新seed
//  采用md5算法生成id
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

var IDGenerator *Generator

func init() {
	IDGenerator = &Generator{}
}
