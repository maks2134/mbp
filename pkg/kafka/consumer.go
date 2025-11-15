package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Consumer struct {
	subscriber message.Subscriber
	logger     watermill.LoggerAdapter
}

func NewConsumer(brokers []string, consumerGroup string, logger watermill.LoggerAdapter) (*Consumer, error) {
	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       brokers,
			ConsumerGroup: consumerGroup,
			Unmarshaler:   kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka subscriber: %w", err)
	}

	return &Consumer{
		subscriber: subscriber,
		logger:     logger,
	}, nil
}

type EventHandler func(ctx context.Context, event interface{}) error

func (c *Consumer) Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	messages, err := c.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	go func() {
		for msg := range messages {
			var event map[string]interface{}
			if err := json.Unmarshal(msg.Payload, &event); err != nil {
				c.logger.Error("Failed to unmarshal message", err, watermill.LogFields{
					"topic": topic,
				})
				msg.Nack()
				continue
			}

			if err := handler(ctx, event); err != nil {
				c.logger.Error("Failed to handle event", err, watermill.LogFields{
					"topic": topic,
				})
				msg.Nack()
				continue
			}

			msg.Ack()
			c.logger.Info("Processed event", watermill.LogFields{
				"topic": topic,
			})
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	return c.subscriber.Close()
}
