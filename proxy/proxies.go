package proxy

import (
	"github.com/kangaloo/go-socks-proxy/monitor"
	"github.com/kangaloo/go-socks-proxy/util"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

var generator = &util.Generator{SeedStr: "proxies"}

type Proxies struct {
	id                string
	src               net.Conn
	dst               net.Conn
	uploadCollector   *monitor.Counter
	downloadCollector *monitor.Counter
	bufSize           int
	dialTimeout       time.Duration
	gw                *sync.WaitGroup
}

func NewProxies(src net.Conn, dst string) (*Proxies, error) {
	p := &Proxies{
		id:          generator.Generate(),
		src:         src,
		bufSize:     1024,
		dialTimeout: time.Second * 10, // will read from config
		gw:          &sync.WaitGroup{},
	}

	dstConn, err := net.DialTimeout("tcp", dst, p.dialTimeout)
	if err != nil {
		return nil, err
	}

	log.WithField("request_domain", dst).Infof("create new proxy destination connection successfully")
	p.dst = dstConn
	p.uploadCollector = monitor.NewFlowCounter(p.src.RemoteAddr().String(), p.dst.RemoteAddr().String(), "upload")
	p.downloadCollector = monitor.NewFlowCounter(p.src.RemoteAddr().String(), p.dst.RemoteAddr().String(), "download")
	return p, nil
}

func (p *Proxies) Run() {
	log.WithField("proxies_id", p.id).Info("proxies start running")
	p.gw.Add(2)
	go forward(p.id, p.gw, p.src, p.dst, p.downloadCollector)
	go forward(p.id, p.gw, p.dst, p.src, p.uploadCollector)
	p.gw.Wait()
	log.WithField("proxies_id", p.id).Info("proxies running completed")
}
