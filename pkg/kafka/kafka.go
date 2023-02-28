package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	turingKafka "go-micro/common/kafka"
	"go-micro/common/logging"
	"go-micro/pkg/config"
)

// kafkaMessage -
type kafkaMessage struct {
	Ts   int64       `json:"ts"`   // 消息推送时间
	Data []ProtoData `json:"data"` // 数据
}

// ProtoData -
type ProtoData struct {
	Ts       int64  `json:"ts"`
	SensorID string `json:"sensorID"` // 传感器ID
	Score    uint16 `json:"score"`
	Result   uint8  `json:"result"`
}

// Client -
type Client struct {
	client *turingKafka.Kafka
	config *config.GoMicroConfig
	logger logging.ILogger
}

// NewKafkaClient -
func NewKafkaClient(logger logging.ILogger, k *turingKafka.Kafka, config *config.GoMicroConfig) *Client {
	return &Client{
		client: k,
		config: config,
		logger: logger,
	}
}

// SendResultToKafka -
func (kc *Client) SendResultToKafka(arr *ProtoBuffer) error {
	if kc.client == nil {
		return fmt.Errorf("kafka is nil")
	}
	strArray := make([]ProtoData, arr.count)
	for i := 0; i < int(arr.count); i++ {
		strArray[i] = ProtoData{
			Ts:       arr.rows[i].Ts,
			SensorID: arr.rows[i].SensorID,
			Score:    arr.rows[i].Score,
			Result:   arr.rows[i].Result,
		}
	}

	msg := &kafkaMessage{
		Ts:   time.Now().UnixMilli(),
		Data: strArray,
	}

	notice, err := json.Marshal(msg)
	if err != nil {
		kc.logger.Errorw("serialization", "err", err)
		return err
	}
	kc.client.ProduceDataWithTimeKey(kc.config.Kafka.Topic, notice)

	return nil
}
