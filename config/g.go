package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
)

func init() {
	getConfig()

}

func getConfig() {
	// name of config file (without extension)
	viper.SetConfigName("config")
	viper.AddConfigPath("/root/GoProjects/go-socks-proxy/conf")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

type Config struct {
}

// Level字段需要自定义解析方法
// DisableColors和FullTimestamp字段，
// 只用在logger.out 为stderr是才生效
type LogConfig struct {
	Writer        io.Writer
	Format        *logrus.TextFormatter
	Level         logrus.Level
	Caller        bool
	DisableColors bool
	FullTimestamp bool
}

func GetLogConfig() *LogConfig {
	config := &LogConfig{
		Format: &logrus.TextFormatter{
			ForceColors: false,
		},
		Level:  logrus.InfoLevel,
		Caller: true,
	}

	// todo 没有log目录时自动创建
	f, err := os.Create("/root/GoProjects/go-socks-proxy/logger/log/socks-proxy.log")
	if err != nil {
		config.Writer = os.Stdout
	}

	w := io.MultiWriter(os.Stdout, f)
	config.Writer = w
	return config
}

/*
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
*/

func (c *Config) Parse() {

}
