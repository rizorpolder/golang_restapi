package config

import (
	"github.com/segmentio/kafka-go"
	"time"
)

type Config struct {
	Brokers []string
	GroupID string
	Topics  []string
}

type ProducerConfig struct {
	Brokers      []string
	BatchSize    int
	BatchTimeout time.Duration
	RequiredAcks kafka.RequiredAcks
}

type ConsumerConfig struct {
	Brokers           []string
	GroupID           string
	Topics            []string
	MinBytes          int64
	MaxBytes          int64
	MaxWaitTime       time.Duration
	ReadBatchTimeout  time.Duration
	HeartbeatInterval time.Duration
	CommitInterval    time.Duration
}

func DefaultProducerConfig(brokers []string) ProducerConfig {
	return ProducerConfig{
		Brokers:      brokers,
		BatchSize:    100,
		BatchTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireOne,
	}
}
func DefaultConsumerConfig(brokers []string, groupID string) ConsumerConfig {
	return ConsumerConfig{}
}
