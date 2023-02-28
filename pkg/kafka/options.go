package kafka

import (
	turingKafka "go-micro/common/kafka"
	"go-micro/common/logging"
	"go-micro/pkg/config"
)

// Option is a function that will set up option.
type Option func(opts *Server)

func loadOptions(options ...Option) *Server {
	opts := &Server{}
	for _, option := range options {
		option(opts)
	}
	if opts.logger == nil {
		opts.logger = new(logging.NoopLogger)
	}
	if opts.config == nil {
		opts.config = config.GetConfig()
	}
	return opts
}

// WithLogger -
func WithLogger(logger logging.ILogger) Option {
	return func(opts *Server) {
		opts.logger = logger
	}
}

// WithConfig -
func WithConfig(c *config.GoMicroConfig) Option {
	return func(opts *Server) {
		opts.config = c
	}
}

// WithKafka -
func WithKafka(k interface{}) Option {
	return func(opts *Server) {
		if k != nil {
			opts.kafka = k.(*turingKafka.Kafka)
		}
	}
}
