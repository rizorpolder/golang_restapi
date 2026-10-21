package config

import (
	"golang_restapi/internal/config"
	"os"
)

type Config struct {
	config.Base
	GRPCHost            string   `env:"GRPC_HOST,required" json:"host"  envDefault:"0.0.0.0"`
	GRPCPort            string   `env:"GRPC_PORT,required" json:"port" envDefault:"50054"`
	DBDSN               string   `env:"DB_DSN,required" json:"db_dnt"`
	AccountGrpcHost     string   `env:"ACCOUNT_GRPC_HOST,required" json:"account_grpc_host"`
	KafkaBrokers        []string `env:"KAFKA_BROKER_HOST,required" envSeparator:","  json:"kafka_broker"`
	KafkaGroupID        string   `env:"KAFKA_CONSUMER_GROUP,required" json:"kafka_consumer_group"`
	KafkaPublishTopic   string   `env:"KAFKA_PUBLISH_TOPIC,required" json:"kafka_publish_topic"`
	KafkaSubscribeTopic string   `env:"KAFKA_SUBSCRIBE_TOPIC,required" json:"kafka_subscribe_topic"`
}

func Load() (*Config, error) {
	path := os.Getenv("ENV_FILE")
	if path == "" {
		path = "configs/transaction.env"
	}
	return config.Load[Config](path)
}

func (c Config) Redacted() Config {
	c.DBDSN = "***"

	return c
}
