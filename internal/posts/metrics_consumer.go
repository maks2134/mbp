package posts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
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
	messagesViewed, err := subscriber.Subscribe(context.Background(), "post.viewed")
	if err != nil {
		return fmt.Errorf("failed to subscribe to post.viewed: %w", err)
	}

	messagesLiked, err := subscriber.Subscribe(context.Background(), "post.liked")
	if err != nil {
		return fmt.Errorf("failed to subscribe to post.liked: %w", err)
	}

	messagesUnliked, err := subscriber.Subscribe(context.Background(), "post.unliked")
	if err != nil {
		return fmt.Errorf("failed to subscribe to post.unliked: %w", err)
	}

	c.logger.Info("Subscriptions created, starting consumers...", nil)

	go func() {
		c.logger.Info("Consumer for post.viewed started, waiting for messages...", nil)
		for msg := range messagesViewed {
			c.processViewedEvent(msg)
		}
	}()

	go func() {
		c.logger.Info("Consumer for post.liked started, waiting for messages...", nil)
		for msg := range messagesLiked {
			c.processLikedEvent(msg)
		}
	}()

	go func() {
		c.logger.Info("Consumer for post.unliked started, waiting for messages...", nil)
		for msg := range messagesUnliked {
			c.processUnlikedEvent(msg)
		}
	}()

	for i := 0; i < 50; i++ {
		runtime.Gosched()
	}
	time.Sleep(1 * time.Second)

	c.logger.Info("All event consumers started - subscriptions should be registered", nil)
	return nil
}

func (c *MetricsSyncConsumer) processViewedEvent(msg *message.Message) {
	var event PostViewedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		c.logger.Error("failed to unmarshal viewed event", err, nil)
		msg.Nack()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.syncViews(ctx, event.PostID, event.Views); err != nil {
		c.logger.Error("failed to sync views", err, nil)
		msg.Nack()
		return
	}

	msg.Ack()
	log.Printf("Synced views for post %d: %d", event.PostID, event.Views)
}

func (c *MetricsSyncConsumer) processLikedEvent(msg *message.Message) {
	var event PostLikedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		c.logger.Error("failed to unmarshal liked event", err, nil)
		msg.Nack()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.syncLikes(ctx, event.PostID, event.Likes); err != nil {
		c.logger.Error("failed to sync likes", err, nil)
		msg.Nack()
		return
	}

	msg.Ack()
	log.Printf("Synced likes for post %d: %d", event.PostID, event.Likes)
}

func (c *MetricsSyncConsumer) processUnlikedEvent(msg *message.Message) {
	var event PostUnlikedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		c.logger.Error("failed to unmarshal unliked event", err, nil)
		msg.Nack()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.syncLikes(ctx, event.PostID, event.Likes); err != nil {
		c.logger.Error("failed to sync likes", err, nil)
		msg.Nack()
		return
	}

	msg.Ack()
	log.Printf("Synced likes for post %d: %d (unliked)", event.PostID, event.Likes)
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
