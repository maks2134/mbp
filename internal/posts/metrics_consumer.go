package posts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MetricsSyncConsumer struct {
	repo   PostsRepositoryInterface
	logger watermill.LoggerAdapter
}

func NewMetricsSyncConsumer(repo PostsRepositoryInterface, logger watermill.LoggerAdapter) *MetricsSyncConsumer {
	return &MetricsSyncConsumer{
		repo:   repo,
		logger: logger,
	}
}

func (c *MetricsSyncConsumer) StartConsumers(subscriber message.Subscriber) error {
	go func() {
		if err := c.consumeViewedEvents(subscriber); err != nil {
			c.logger.Error("failed to consume viewed events", err, nil)
		}
	}()

	go func() {
		if err := c.consumeLikedEvents(subscriber); err != nil {
			c.logger.Error("failed to consume liked events", err, nil)
		}
	}()

	go func() {
		if err := c.consumeUnlikedEvents(subscriber); err != nil {
			c.logger.Error("failed to consume unliked events", err, nil)
		}
	}()

	return nil
}

func (c *MetricsSyncConsumer) consumeViewedEvents(subscriber message.Subscriber) error {
	messages, err := subscriber.Subscribe(context.Background(), "post.viewed")
	if err != nil {
		return fmt.Errorf("failed to subscribe to post.viewed: %w", err)
	}

	for msg := range messages {
		var event PostViewedEvent
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			c.logger.Error("failed to unmarshal viewed event", err, nil)
			msg.Nack()
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := c.syncViews(ctx, event.PostID, event.Views); err != nil {
			c.logger.Error("failed to sync views", err, nil)
			msg.Nack()
			cancel()
			continue
		}
		cancel()

		msg.Ack()
		log.Printf("Synced views for post %d: %d", event.PostID, event.Views)
	}

	return nil
}

func (c *MetricsSyncConsumer) consumeLikedEvents(subscriber message.Subscriber) error {
	messages, err := subscriber.Subscribe(context.Background(), "post.liked")
	if err != nil {
		return fmt.Errorf("failed to subscribe to post.liked: %w", err)
	}

	for msg := range messages {
		var event PostLikedEvent
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			c.logger.Error("failed to unmarshal liked event", err, nil)
			msg.Nack()
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := c.syncLikes(ctx, event.PostID, event.Likes); err != nil {
			c.logger.Error("failed to sync likes", err, nil)
			msg.Nack()
			cancel()
			continue
		}
		cancel()

		msg.Ack()
		log.Printf("Synced likes for post %d: %d", event.PostID, event.Likes)
	}

	return nil
}

func (c *MetricsSyncConsumer) consumeUnlikedEvents(subscriber message.Subscriber) error {
	messages, err := subscriber.Subscribe(context.Background(), "post.unliked")
	if err != nil {
		return fmt.Errorf("failed to subscribe to post.unliked: %w", err)
	}

	for msg := range messages {
		var event PostUnlikedEvent
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			c.logger.Error("failed to unmarshal unliked event", err, nil)
			msg.Nack()
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := c.syncLikes(ctx, event.PostID, event.Likes); err != nil {
			c.logger.Error("failed to sync likes", err, nil)
			msg.Nack()
			cancel()
			continue
		}
		cancel()

		msg.Ack()
		log.Printf("Synced likes for post %d: %d (unliked)", event.PostID, event.Likes)
	}

	return nil
}

func (c *MetricsSyncConsumer) syncViews(ctx context.Context, postID, views int) error {
	post, err := c.repo.FindByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	post.CountViewers = views
	post.UpdatedAt = time.Now()

	if err := c.repo.Update(ctx, post); err != nil {
		return fmt.Errorf("failed to update views: %w", err)
	}

	return nil
}

func (c *MetricsSyncConsumer) syncLikes(ctx context.Context, postID, likes int) error {
	post, err := c.repo.FindByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	post.Like = likes
	post.UpdatedAt = time.Now()

	if err := c.repo.Update(ctx, post); err != nil {
		return fmt.Errorf("failed to update likes: %w", err)
	}

	return nil
}
