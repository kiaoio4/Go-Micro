package component

import (
	"context"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/spf13/viper"

	platformConf "go-micro/common/conf"
	"go-micro/common/logging"
	"go-micro/common/logging/conf"
	"go-micro/common/micro"
	microConf "go-micro/common/micro/conf"
	"go-micro/common/utils"
)

// LoggerGroupComponent is Component for logging
type LoggerGroupComponent struct {
	micro.EmptyComponent
	group  *logging.LoggerGroup
	enable bool
}

// Name of the component
func (c *LoggerGroupComponent) Name() string {
	return "LoggerGroup"
}

// PreInit called before Init()
func (c *LoggerGroupComponent) PreInit(ctx context.Context) error {
	// load config
	conf.SetDefaultLogGroupConfig()
	return nil
}

// SetDynamicConfig called when get dynamic config for the first time
func (c *LoggerGroupComponent) SetDynamicConfig(config *platformConf.NodeConfig) error {
	c.enable = config.APM != nil && config.APM.EnableLog
	return nil
}

// Init the component
func (c *LoggerGroupComponent) Init(server *micro.Server) error {
	// init
	var err error
	// setup logger
	basicConf := microConf.GetBasicConfig()
	// spew.Dump(basicConf)
	logConf, err := conf.GetLogGroupConfig()
	if err != nil {
		return err
	}

	if !c.enable {
		logConf.GraylogAddr = ""
		logConf.LokiAddr = ""
	}
	str := viper.GetString("LOG_LEVELS_DEBUG")
	if str != "" {
		logConf.Levels[logging.LevelDebug] = strings.Split(str, ",")
	}
	str = viper.GetString("LOG_LEVELS_INFO")
	if str != "" {
		logConf.Levels[logging.LevelInfo] = strings.Split(str, ",")
	}
	str = viper.GetString("LOG_LEVELS_WARN")
	if str != "" {
		logConf.Levels[logging.LevelWarn] = strings.Split(str, ",")
	}
	str = viper.GetString("LOG_LEVELS_ERROR")
	if str != "" {
		logConf.Levels[logging.LevelError] = strings.Split(str, ",")
	}

	// spew.Dump(logConf)
	c.group, err = logging.CreateLoggerGroup(basicConf, logConf)
	if err != nil {
		return err
	}
	server.RegisterElement(&micro.LoggerGroupElementKey, c.group)

	if basicConf.IsDevMode {
		spew.Dump(logConf)
	}
	return nil
}

const (
	urlGroupLog = "Log??????(????????????)"
	urlLogLevel = "/log/level"
)

// SetupHandler of echo if the component need
func (c *LoggerGroupComponent) SetupHandler(root echoswagger.ApiRoot, base string) error {
	g := root.Group(urlGroupLog, urlLogLevel)
	g.POST("", c.setHandler).
		AddParamQuery("", "module", "module", true).
		AddParamQuery("", "level", "level", true).
		AddResponse(http.StatusOK, "successful operation", "", nil).
		SetOperationId("add").
		SetSummary("set log level for given module")
	g.GET("", c.getHandler).
		AddResponse(http.StatusOK, "successful operation", "", nil).
		SetOperationId("get").
		SetSummary("get log levels for all modules")
	return nil
}

func (c *LoggerGroupComponent) setHandler(ctx echo.Context) error {
	var module, level string
	err := echo.QueryParamsBinder(ctx).
		String("module", &module).
		String("level", &level).
		BindError()
	if err != nil {
		return utils.GetJSONResponse(ctx, err, nil)
	}
	c.group.SetLevel(module, level)
	return utils.GetJSONResponse(ctx, nil, "OK")
}

func (c *LoggerGroupComponent) getHandler(ctx echo.Context) error {
	levels := c.group.GetLevels()
	return utils.GetJSONResponse(ctx, nil, levels)
}

// PostStop called after Stop()
func (c *LoggerGroupComponent) PostStop(ctx context.Context) error {
	// post stop
	return c.group.Sync()
}
