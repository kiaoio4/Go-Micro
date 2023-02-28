package metric

// FileReadMetric 读取文件度量指标
type FileReadMetric interface {
	Inc(args ...string) // 自增读取文件
	Register() error    // 注册
}