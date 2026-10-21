package config

import (
	"golang_restapi/internal/config"
	"os"
)

type Config struct {
	config.Base
	DBDSN               string   `env:"DB_DSN,required"`
	GRPCHost            string   `env:"GRPC_HOST" envDefault:"0.0.0.0"`
	GRPCPort            string   `env:"GRPC_PORT" envDefault:"50053"`
	KafkaBrokers        []string `env:"KAFKA_BROKER_HOST,required" envSeparator:","  json:"kafka_broker"`
	KafkaGroupID        string   `env:"KAFKA_CONSUMER_GROUP,required" json:"kafka_consumer_group"`
	KafkaPublishTopic   string   `env:"KAFKA_PUBLISH_TOPIC,required" json:"kafka_publish_topic"`
	KafkaSubscribeTopic string   `env:"KAFKA_SUBSCRIBE_TOPIC,required" json:"kafka_subscribe_topic"`
}

func Load() (*Config, error) {
	path := os.Getenv("ENV_FILE")
	if path == "" {
		path = "configs/account.env"
	}
	return config.Load[Config](path)
}

func (c Config) Redacted() Config {
	c.DBDSN = "***"

	return c
}
