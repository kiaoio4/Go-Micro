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

	"github.com/pangpanglabs/echoswagger/v2"
)

// TestElementKey is Element Key for TestComponent
var TestElementKey = micro.ElementKey("TestComponent")

// TestComponent is Component for TestComponent
type TestComponent struct {
	micro.EmptyComponent
	stopChan      chan struct{}
	handler       *pkg.GoMicro
	logger        logging.ILogger
	nacosClient   *configuration.NacosClient
	gossipKVCache *microComponent.GossipKVCacheComponent
	cluster       string
}

// Name of the component
func (c *TestComponent) Name() string {
	return "TestComponent"
}

// PreInit called before Init()
func (c *TestComponent) PreInit(ctx context.Context) error {
	// load config
	config.SetDefaultTestConfig()
	return nil
}

// Init the component
func (c *TestComponent) Init(server *micro.Server) error {
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

	c.handler, err = pkg.NewGoMicro(basicConfig, c.logger, c.gossipKVCache)
	if err != nil {
		return err
	}

	return nil
}

// OnConfigChanged called when dynamic config changed
func (c *TestComponent) OnConfigChanged(*platformConf.NodeConfig) error {
	return micro.ErrNeedRestart
}

// SetupHandler of echo if the component need
func (c *TestComponent) SetupHandler(root echoswagger.ApiRoot, base string) error {
	basicConf := microConf.GetBasicConfig()
	selfServiceName := basicConf.Service
	c.handler.SetupWeb(root, base, selfServiceName)

	return nil
}

// Start the component
func (c *TestComponent) Start(ctx context.Context) error {
	// start
	go c.handler.Start(c.stopChan)
	return nil
}

// Stop the component
func (c *TestComponent) Stop(ctx context.Context) error {
	// stop
	c.stopChan <- struct{}{}
	return nil
}
