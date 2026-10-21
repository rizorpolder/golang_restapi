package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"golang_restapi/internal/kafka/config"
	"time"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	logger *zerolog.Logger
}

type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Headers   map[string]string
	Partition int
}

func NewProducer(cfg config.ProducerConfig, logger *zerolog.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Balancer:               &kafka.LeastBytes{},
		BatchSize:              cfg.BatchSize,
		BatchTimeout:           cfg.BatchTimeout,
		RequiredAcks:           cfg.RequiredAcks,
		Async:                  false,
		AllowAutoTopicCreation: true,
	}
	return &Producer{
		writer: writer,
		logger: logger,
	}
}

func (p *Producer) SendMessage(ctx context.Context, message Message) error {
	kafkaMsg := kafka.Message{
		Topic:     message.Topic,
		Key:       message.Key,
		Value:     message.Value,
		Partition: message.Partition,
		Time:      time.Now(),
	}

	for key, value := range message.Headers {
		header := kafka.Header{
			Key:   key,
			Value: []byte(value),
		}
		kafkaMsg.Headers = append(kafkaMsg.Headers, header)
	}

	err := p.writer.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		p.logger.Error().
			Err(err).
			Str("topic", message.Topic).
			Str("key", string(message.Key)).
			Msg("failed to send message to kafka")
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}
	p.logger.Debug().
		Str("topic", message.Topic).
		Str("key", string(message.Key)).
		Msg("message sent to kafka successfully")

	return nil
}
func (p *Producer) SendJSONMessage(ctx context.Context, topic string, key string, data interface{}) error {
	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	msg := Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
		Headers: map[string]string{
			"content-type": "application/json",
		},
	}
	return p.SendMessage(ctx, msg)
}
func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
