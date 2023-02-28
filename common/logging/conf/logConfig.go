package conf

import (
	"github.com/spf13/viper"
)

const (
	logLevel       = "log.level"
	logPath        = "log.path"
	logGraylogAddr = "log.graylog"
	logLokiAddr    = "log.loki"
	logConsole     = "log.console"
)

var defaultLogConfig = LogConfig{
	Level:   "INFO",
	Path:    "",
	Console: true,
}

//LogConfig  日志配置
type LogConfig struct {
	Level       string `toml:"level"`
	Path        string `toml:"path"`
	GraylogAddr string `toml:"graylog"`
	LokiAddr    string `toml:"loki"`
	Console     bool   `toml:"console"`
}

//SetDefaultLogConfig 获取默认日志配置
func SetDefaultLogConfig() {
	viper.SetDefault(logLevel, defaultLogConfig.Level)
	viper.SetDefault(logPath, defaultLogConfig.Path)
	viper.SetDefault(logConsole, defaultLogConfig.Console)
}

//GetLogConfig  获取日志配置
func GetLogConfig() *LogConfig {
	return &LogConfig{
		Level:       viper.GetString(logLevel),
		Path:        viper.GetString(logPath),
		GraylogAddr: viper.GetString(logGraylogAddr),
		LokiAddr:    viper.GetString(logLokiAddr),
		Console:     viper.GetBool(logConsole),
	}
}
