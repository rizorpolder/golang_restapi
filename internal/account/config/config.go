package config

import (
	"golang_restapi/internal/config"
	"os"
)

type Config struct {
	config.Base
	DBDSN                 string   `env:"DB_DSN,required"`
	GRPCHost              string   `env:"GRPC_HOST" envDefault:"0.0.0.0"`
	GRPCPort              string   `env:"GRPC_PORT" envDefault:"50052"`
	KafkaBrokers          []string `env:"KAFKA_BROKER_HOST" envSeparator:","  json:"kafka_broker" required:"true"`
	KafkaGroupID          string   `env:"KAFKA_CONSUMER_GROUP" json:"kafka_consumer_group" required:"true"`
	KafkaTransactionTopic string   `env:"KAFKA_TRANSACTION_TOPIC" json:"kafka_transaction_topic" required:"true"`
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
