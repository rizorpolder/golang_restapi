package kafka

import (
	"context"
	"fmt"
	"golang_restapi/internal/kafka/config"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

type Client interface {
	Publish(ctx context.Context, topic string, key string, data interface{}) error
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Close() error
}

type MessageHandler func(ctx context.Context, topic string, key string, data []byte) error

type Kafka struct {
	producer *Producer
	brokers  []string
	groupId  string
	logger   *zerolog.Logger
}

func New(producer *Producer, brokers []string, groupId string, logger *zerolog.Logger) *Kafka {
	return &Kafka{
		producer: producer,
		brokers:  brokers,
		groupId:  groupId,
		logger:   logger,
	}
}

func (k *Kafka) Publish(ctx context.Context, topic string, key string, data interface{}) error {
	return k.producer.SendJSONMessage(ctx, topic, key, data)
}
func (k *Kafka) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	uniqGroupID := k.groupId + "_" + topic
	cfg := config.DefaultConsumerConfig(k.brokers, uniqGroupID, []string{topic})
	consumer := NewConsumer(cfg, k.logger)

	go func() {

		defer func() {
			if err := consumer.Close(); err != nil {
				k.logger.Error().Err(err).Str("topic", topic).Msg("failed to close kafka consumer")
			}
		}()

		k.logger.Info().
			Str("topic", topic).
			Str("groupId", uniqGroupID).
			Msg("subscribing to topic")

		messageHandler := ConsumerMessageHandler(func(ctx context.Context, msg kafka.Message) error {
			return handler(ctx, msg.Topic, string(msg.Key), msg.Value)
		})

		_ = consumer.Start(ctx, messageHandler)

	}()
	return nil
}
func (k *Kafka) Close() error {
	if k.producer != nil {
		if err := k.producer.Close(); err != nil {
			return fmt.Errorf("close kafka producer failed: %v", err)
		}
	}
	return nil
}
