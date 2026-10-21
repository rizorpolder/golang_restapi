package kafka

import (
	"context"
	"golang_restapi/internal/kafka/config"
	"time"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

type ConsumerMessageHandler func(ctx context.Context, message kafka.Message) error
type Consumer struct {
	reader *kafka.Reader
	logger *zerolog.Logger
}

func (c *Consumer) GetConfig() kafka.ReaderConfig {
	return c.reader.Config()
}

func NewConsumer(cfg config.ConsumerConfig, logger *zerolog.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           cfg.Brokers,
		GroupID:           cfg.GroupID,
		Topic:             cfg.Topics[0],
		MinBytes:          cfg.MinBytes,
		MaxBytes:          cfg.MaxBytes,
		MaxWait:           cfg.MaxWaitTime,
		ReadBatchTimeout:  cfg.ReadBatchTimeout,
		HeartbeatInterval: cfg.HeartbeatInterval,
		CommitInterval:    cfg.CommitInterval,
	})
	return &Consumer{
		reader: reader,
		logger: logger,
	}
}

func (c *Consumer) Start(ctx context.Context, handler ConsumerMessageHandler) error {
	c.logger.Info().
		Str("topic", c.reader.Config().Topic).
		Str("group_id", c.reader.Config().GroupID).
		Msg("starting kafka consumer")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info().Msg("kafka consumer stopped")
			return ctx.Err()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				c.logger.Error().Err(err).Msg("failed to read kafka message")
				time.Sleep(time.Second)
				continue
			}
			if err := handler(ctx, msg); err != nil {
				c.logger.Error().Err(err).
					Str("topic", msg.Topic).
					Int("partition", msg.Partition).
					Int64("offset", msg.Offset).
					Msg("failed to process message")
			} else {
				c.logger.Debug().
					Str("topic", msg.Topic).
					Int("partition", msg.Partition).
					Int64("offset", msg.Offset).
					Msg("message successfully processed")

			}
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.logger.Error().Err(err).Msg("failed to commit kafka message")
			}
		}
	}
}

func (c *Consumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
