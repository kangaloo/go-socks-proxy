package logger

import (
	"github.com/kangaloo/go-socks-proxy/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	conf := config.GetLogConfig()
	log.SetLevel(conf.Level)
	log.SetFormatter(conf.Format)
	log.SetReportCaller(conf.Caller)
	log.SetOutput(conf.Writer)
}

/*
todo 实现一个 logFormatter
type Formatter interface {
	Format(*Entry) ([]byte, error)
}
*/

// 根据配置文件，决定使用logrus还是标准库log

func init() {

}
