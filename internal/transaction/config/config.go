package config

import (
	"golang_restapi/internal/config"
	"os"
)

type Config struct {
	config.Base
	GRPCHost              string   `env:"GRPC_HOST" json:"host" required:"true" default:"localhost"`
	GRPCPort              string   `env:"GRPC_PORT" json:"port" required:"true" default:"9000"`
	DBDSN                 string   `env:"DB_DSN" json:"db_dnt" required:"true"`
	AccountGrpcHost       string   `env:"ACCOUNT_GRPC_HOST" json:"account_grpc_host" required:"true"`
	KafkaBrokers          []string `env:"KAFKA_BROKERS_HOST" envSeparator:","  json:"kafka_brokers" required:"true"`
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
