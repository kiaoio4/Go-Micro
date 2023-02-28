package config

import (
	"github.com/spf13/viper"
)

const (
	configPath = "micro.path"
)

var defaultTestConfig = TestConfig{
	Path: "/home/workspace",
}

// WorkConfig 配置
type TestConfig struct {
	Path string `toml:"path"`
}

// SetDefaultWorkConfig -
func SetDefaultTestConfig() {
	viper.SetDefault(configPath, defaultTestConfig.Path)
}

// GetWorkConfig Get默认配置参数
func GetTestConfig() *TestConfig {
	return &TestConfig{
		Path: viper.GetString(configPath),
	}
}
