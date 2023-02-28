package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-micro/common/kafka"
	"go-micro/common/logging"
	"go-micro/pkg/config"

	"github.com/davecgh/go-spew/spew"
)

// Handler -
type Handler interface {
	Start(context.Context)
	Write(uint64, []byte) error
	Stop()
}

// Server - 时序数据库管理结构
type Server struct {
	kafka           *kafka.Kafka
	srv             *Client
	numericalBuffer *sync.Map
	status          bool
	config          *config.GoMicroConfig
	logger          logging.ILogger
}

// New -
func New(opts ...Option) (Handler, error) {
	srv := loadOptions(opts...)

	if srv.kafka == nil {
		return nil, fmt.Errorf("new Clinet fail")
	}

	spew.Dump(kafka.GetConfig())
	spew.Dump(srv.config)

	srv.numericalBuffer = new(sync.Map)
	srv.status = true
	srv.srv = NewKafkaClient(srv.logger, srv.kafka, srv.config)

	return srv, nil
}

// Stop -
func (s *Server) Stop() {
	s.numericalBuffer.Range(func(key, value interface{}) bool {
		if err := value.(*ProtoBuffer).flush(); err != nil {
			s.logger.Errorw("kafka numerical flush", "err", err)
		}
		s.numericalBuffer.Delete(key)
		return true
	})

	s.srv.client.Close()
	s.logger.Infow("kafka service close")
}

// Start -
func (s *Server) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(s.config.Kafka.Interval))
	defer func() {
		ticker.Stop()
	}()

	s.logger.Infow("kafka service start")

	for {
		select {
		case <-ticker.C:
		case <-ctx.Done():
			s.logger.Infow("kafka service stop")
			return
		}
		s.numericalBuffer.Range(func(key, value interface{}) bool {
			b := value.(*ProtoBuffer)
			if err := b.flush(); err != nil {
				s.logger.Errorw("kafka nemerical flush", "err", err)
			}
			return true
		})
	}
}

// Write -
func (s *Server) Write(id uint64, frame []byte) error {
	// 数值表数据入库
	if s.config.Kafka.Enable {
		if err := s.sendToKafka(id, 1, time.Now().Unix()); err != nil {
			return err
		}
	}
	return nil
}
