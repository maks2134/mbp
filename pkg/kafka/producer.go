package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Producer struct {
	publisher message.Publisher
	logger    watermill.LoggerAdapter
}

func NewProducer(brokers []string, logger watermill.LoggerAdapter) (*Producer, error) {
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   brokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka publisher: %w", err)
	}

	return &Producer{
		publisher: publisher,
		logger:    logger,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	if err := p.publisher.Publish(topic, msg); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	p.logger.Info("Published event", watermill.LogFields{
		"topic": topic,
		"event": string(payload),
	})

	return nil
}

func (p *Producer) Close() error {
	return p.publisher.Close()
}
