package logger

import log "github.com/sirupsen/logrus"

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors: false,
	})

	log.SetReportCaller(true)
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
