package component

import (
	"context"
	"io"

	"github.com/davecgh/go-spew/spew"
	"github.com/opentracing/opentracing-go"

	platformConf "go-micro/common/conf"
	"go-micro/common/micro"
	microConf "go-micro/common/micro/conf"
	"go-micro/common/tracing"
)

// TracingComponent is Component for tracing
type TracingComponent struct {
	micro.EmptyComponent
	closer io.Closer
	enable bool
}

// Name of the component
func (c *TracingComponent) Name() string {
	return "Trace"
}

// PreInit called before Init()
func (c *TracingComponent) PreInit(ctx context.Context) error {
	// load config
	tracing.SetDefaultTraceConfig()
	return nil
}

// SetDynamicConfig called when get dynamic config for the first time
func (c *TracingComponent) SetDynamicConfig(config *platformConf.NodeConfig) error {
	c.enable = config.APM != nil && !config.APM.EnableTrace
	return nil
}

// Init the component
func (c *TracingComponent) Init(server *micro.Server) error {
	// init
	var err error
	var tracer opentracing.Tracer
	// setup tracer
	basicConf := microConf.GetBasicConfig()
	traceConf := tracing.GetTraceConfig()
	//logger := server.Get(&logging.Key).(logging.ILogger)

	if c.enable {
		server.RegisterElement(&micro.TracingElementKey, opentracing.NoopTracer{})
		return nil
	}

	tracer, c.closer, err = tracing.CreateTracer(*basicConf, *traceConf, nil)
	if err != nil {
		panic("Could not initialize jaeger tracer: " + err.Error())
	}
	server.RegisterElement(&micro.TracingElementKey, tracer)

	if basicConf.IsDevMode {
		spew.Dump(traceConf)
	}
	return nil
}

// PostStop called after Stop()
func (c *TracingComponent) PostStop(ctx context.Context) error {
	// post stop
	return c.closer.Close()
}
