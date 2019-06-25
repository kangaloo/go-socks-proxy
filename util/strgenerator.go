package util

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"sync"
	"time"
)

type Generator struct {
	sync.Mutex
	seed    int
	SeedStr string
}

func (g *Generator) Generate() string {
	g.Lock()
	defer g.Unlock()
	g.seed += 1
	h := md5.New()
	h.Write([]byte(strconv.Itoa(int(time.Now().UnixNano()))))
	id := hex.EncodeToString(h.Sum([]byte(strconv.Itoa(g.seed) + g.SeedStr)))
	return id
}
