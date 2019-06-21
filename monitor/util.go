package monitor

import (
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
	id := strconv.Itoa(int(time.Now().UnixNano()))
	return id
}

var IDGenerator *Generator

func init() {
	IDGenerator = &Generator{}
}
