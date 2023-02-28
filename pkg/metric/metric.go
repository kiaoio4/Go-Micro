package metric

import (
	"github.com/pkg/errors"
)

const (
	// MonitorSuccess 统计成功项
	MonitorSuccess = "success"
	// MonitorFailed 统计失败项
	MonitorFailed = "failed"
)

const (
	// MonitorNamespace .
	MonitorNamespace = "Go"
	// MonitorSubsystem .
	MonitorSubsystem = "micro"

	// MonitorFileRead 文件读取
	MonitorFileRead = "file_read"
)

// FileCacheMonitor 文件缓存监控
type FileCacheMonitor struct {
	fileReadMetric FileReadMetric // 文件读取指标
}

// NewFileCacheMonitor .
func NewFileCacheMonitor(fr FileReadMetric) (*FileCacheMonitor, error) {
	f := &FileCacheMonitor{
		fileReadMetric: fr,
	}
	if err := f.registerFileCacheMonitor(); err != nil {
		return nil, errors.Wrap(err, "注册文件缓存监控服务")
	}
	return f, nil
}

// registerFileCacheMonitor 注册文件缓存监控
func (f *FileCacheMonitor) registerFileCacheMonitor() error {
	return f.fileReadMetric.Register()
}

// SetFileReadValues 设置文件读取计数
func (f *FileCacheMonitor) SetFileReadValues(sensorID, result string) {
	f.fileReadMetric.Inc(sensorID, result)
}
