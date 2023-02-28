package config

import (
	logging "go-micro/common/logging/conf"
	basic "go-micro/common/micro/conf"
	tracing "go-micro/common/tracing"
)

// GoMicroConfig 配置
type GoMicroConfig struct {
	Basic *basic.BasicConfig
	Test  *TestConfig          `toml:"-"`
	Log   *logging.LogConfig   `toml:"-"`
	Trace *tracing.TraceConfig `toml:"-"`
	Pprof *PprofConfig         `toml:"-"`
}

// SetDefaultGoMicroTestConfig -
func SetDefaultGoMicroTestConfig() {
	basic.SetDefaultBasicConfig()
	SetDefaultTestConfig()
	SetDefaultPprofConfig()
}

// GetConfig Get默认配置参数
func GetConfig() *GoMicroConfig {
	return &GoMicroConfig{
		Basic: basic.GetBasicConfig(),
		Test:  GetTestConfig(),
		Pprof: GetPprofConfig(),
		Log:   logging.GetLogConfig(),
		Trace: tracing.GetTraceConfig(),
	}
}
