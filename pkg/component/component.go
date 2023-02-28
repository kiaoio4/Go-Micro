package component

import (
	"context"

	platformConf "go-micro/common/conf"
	"go-micro/common/configuration"
	"go-micro/common/logging"
	"go-micro/common/micro"
	microComponent "go-micro/common/micro/component"
	microConf "go-micro/common/micro/conf"
	"go-micro/pkg"
	"go-micro/pkg/config"
	"go-micro/pkg/kafka"

	"github.com/pangpanglabs/echoswagger/v2"
)

// DataRawdbElementKey is Element Key for Datarawdb
var DataRawdbElementKey = micro.ElementKey("TestComponent")

// DataRawdbComponent is Component for DataRawdb
type DataRawdbComponent struct {
	micro.EmptyComponent
	stopChan      chan struct{}
	handler       *pkg.GoMicro
	logger        logging.ILogger
	nacosClient   *configuration.NacosClient
	gossipKVCache *microComponent.GossipKVCacheComponent
	kafka         kafka.Handler
	cluster       string
}

// Name of the component
func (c *DataRawdbComponent) Name() string {
	return "TestComponent"
}

// PreInit called before Init()
func (c *DataRawdbComponent) PreInit(ctx context.Context) error {
	// load config
	config.SetDefaultTestConfig()
	return nil
}

// Init the component
func (c *DataRawdbComponent) Init(server *micro.Server) error {
	// init
	basicConf := microConf.GetBasicConfig()
	c.cluster = server.PrivateCluster
	if basicConf.IsDynamicConfig {
		c.nacosClient = server.GetElement(&micro.NacosClientElementKey).(*configuration.NacosClient)
	}
	elkvcache := server.GetElement(&micro.GossipKVCacheElementKey)
	if elkvcache != nil {
		c.gossipKVCache = elkvcache.(*microComponent.GossipKVCacheComponent)
	}
	logger := server.GetElement(&micro.LoggingElementKey).(logging.ILogger)
	c.logger = logger

	basicConfig := config.GetConfig()
	var err error
	c.stopChan = make(chan struct{})
	// New Kafka
	elkafka := server.GetElement(&micro.KafkaElementKey)
	if elkafka != nil {
		if c.kafka, err = kafka.New(
			kafka.WithKafka(elkafka),
			kafka.WithLogger(c.logger),
		); err != nil {
			c.logger.Errorw("NewServer", "err", err)
			return err
		}
	} else {
		c.logger.Warnw("GetElement", "GetElement", "is null")
	}

	c.handler, err = pkg.NewGoMicro(basicConfig, c.logger, c.gossipKVCache, c.kafka)
	if err != nil {
		return err
	}

	return nil
}

// OnConfigChanged called when dynamic config changed
func (c *DataRawdbComponent) OnConfigChanged(*platformConf.NodeConfig) error {
	return micro.ErrNeedRestart
}

// SetupHandler of echo if the component need
func (c *DataRawdbComponent) SetupHandler(root echoswagger.ApiRoot, base string) error {
	basicConf := microConf.GetBasicConfig()
	selfServiceName := basicConf.Service
	c.handler.SetupWeb(root, base, selfServiceName)

	return nil
}

// Start the component
func (c *DataRawdbComponent) Start(ctx context.Context) error {
	if c.kafka != nil {
		go c.kafka.Start(ctx)
	}

	// start
	go c.handler.Start(c.stopChan)
	return nil
}

// Stop the component
func (c *DataRawdbComponent) Stop(ctx context.Context) error {
	if c.kafka != nil {
		c.kafka.Stop()
	}

	// stop
	c.stopChan <- struct{}{}
	return nil
}
