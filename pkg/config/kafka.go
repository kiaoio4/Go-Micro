package config

import "github.com/spf13/viper"

const (
	configKafkaEnable          = "kafka.enable"
	configKafkaServer          = "kafka.bootstrap"
	configGroupID              = "kafka.groupid"
	configKafkaMessageMaxBytes = "kafka.messagemaxbytes"
	configKafkaTopic           = "kafka.topic"
	configKafkaInterval        = "kafka.interval"
)

var defaultKafkaConfig = KafkaConfig{
	Enable:          false,
	Server:          "localhost:9092",
	GroupID:         "go-micro",
	MessageMaxBytes: 67108864,
	Topic:           "func",
	Interval:        2,
}

// KafkaConfig -
type KafkaConfig struct {
	Enable          bool   `toml:"enable"`
	Server          string `toml:"server"`
	GroupID         string `toml:"groupid"`
	MessageMaxBytes int    `toml:"messageMaxBytes"`
	Topic           string `toml:"topic"`
	Interval        int    `toml:"interval"`
}

// SetDefaultKafkaConfig -
func SetDefaultKafkaConfig() {
	viper.SetDefault(configKafkaEnable, defaultKafkaConfig.Enable)
	viper.SetDefault(configKafkaServer, defaultKafkaConfig.Server)
	viper.SetDefault(configGroupID, defaultKafkaConfig.GroupID)
	viper.SetDefault(configKafkaMessageMaxBytes, defaultKafkaConfig.MessageMaxBytes)
	viper.SetDefault(configKafkaTopic, defaultKafkaConfig.Topic)
	viper.SetDefault(configKafkaInterval, defaultKafkaConfig.Interval)
}

// GetKafkaConfig -
func GetKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Enable:          viper.GetBool(configKafkaEnable),
		Server:          viper.GetString(configKafkaServer),
		GroupID:         viper.GetString(configGroupID),
		MessageMaxBytes: viper.GetInt(configKafkaMessageMaxBytes),
		Topic:           viper.GetString(configKafkaTopic),
		Interval:        viper.GetInt(configKafkaInterval),
	}
}
