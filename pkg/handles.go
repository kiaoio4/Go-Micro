package pkg

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go-micro/common/logging"
	microComponent "go-micro/common/micro/component"
	"go-micro/common/utils"
	"go-micro/pkg/code"
	"go-micro/pkg/config"
	"go-micro/pkg/kafka"
	"go-micro/pkg/metric"
	"go-micro/pkg/metric/monitor"

	"github.com/labstack/echo"
)

const (
	// API go micro Api.
	API = "Go Micro API"
)

// AppVersion add to file suffix name
var AppVersion = "v1.0"

// GoMicro micro struct
type GoMicro struct {
	logger        logging.ILogger
	working       *sync.Mutex
	gossipKVCache *microComponent.GossipKVCacheComponent
	kafka         kafka.Handler
	config        *config.GoMicroConfig
	once          sync.Once
	exportMetrics *metric.FileCacheMonitor
}

// NewGoMicro Instantiation object
func NewGoMicro(config *config.GoMicroConfig, logger logging.ILogger, gossipKVCache *microComponent.GossipKVCacheComponent, k kafka.Handler) (*GoMicro, error) {
	// 初始化指标采集模块
	fr := monitor.NewFileRead("type")
	m, err := metric.NewFileCacheMonitor(fr)
	if err != nil {
		return nil, err
	}
	db := &GoMicro{
		config:        config,
		working:       new(sync.Mutex),
		logger:        logger,
		exportMetrics: m,
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// gossipKVCache Timer
	if gossipKVCache != nil {
		db.gossipKVCache = gossipKVCache
	}

	// kafka
	if db.config.Kafka.Enable {
		logger.Info("NewGoMicro Kafka.Enable", k)
		db.kafka = k
	}

	return db, nil
}

// Start Connect kafka Store & taoClient
func (micro *GoMicro) Start(stop chan struct{}) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	//decode

	go micro.handleDecodeResult(sigchan, 1)

	<-stop
	micro.Close()
}

// Close kafka & haystack store
func (micro *GoMicro) Close() {

	micro.once.Do(func() {
		micro.logger.Warnw("close service")
	})

	micro.logger.Info("Quit!")
}

// handleDecodeResult -
func (micro *GoMicro) handleDecodeResult(sigchan chan os.Signal, workindex int) {
	for {
		select {
		case sig := <-sigchan:
			micro.Close()
			micro.logger.Errorf("Caught signal %v: terminating\n", sig)
		}
	}
}

func (micro *GoMicro) getTestData(c echo.Context) error {
	path := c.QueryParam("path")
	return c.JSON(http.StatusOK, utils.ResponseV2{
		Code: code.Success,
		Msg:  "OK",
		Data: path},
	)

}
