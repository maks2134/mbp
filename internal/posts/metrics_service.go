package posts

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

const (
	keyPostLikes     = "post:likes:%d"
	keyPostViews     = "post:views:%d"
	keyUserLikedPost = "user:%d:liked:%d"
)

type MetricsService struct {
	redis     *redis.Client
	publisher message.Publisher
	logger    watermill.LoggerAdapter
}

func NewMetricsService(redisClient *redis.Client, publisher message.Publisher, logger watermill.LoggerAdapter) *MetricsService {
	return &MetricsService{
		redis:     redisClient,
		publisher: publisher,
		logger:    logger,
	}
}

// IncrementViews увеличивает счетчик просмотров поста
func (s *MetricsService) IncrementViews(ctx context.Context, postID int) error {
	key := fmt.Sprintf(keyPostViews, postID)
	views, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}

	// Публикуем событие просмотра
	event := PostViewedEvent{
		PostID: postID,
		Views:  int(views),
	}
	if err := s.publishEvent("post.viewed", event); err != nil {
		s.logger.Error("failed to publish post.viewed event", err, nil)
	}

	return nil
}

// GetViews возвращает количество просмотров поста
func (s *MetricsService) GetViews(ctx context.Context, postID int) (int, error) {
	key := fmt.Sprintf(keyPostViews, postID)
	views, err := s.redis.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get views: %w", err)
	}
	return views, nil
}

// LikePost ставит лайк посту от пользователя
func (s *MetricsService) LikePost(ctx context.Context, userID, postID int) error {
	userLikedKey := fmt.Sprintf(keyUserLikedPost, userID, postID)
	exists, err := s.redis.Exists(ctx, userLikedKey).Result()
	if err != nil {
		return fmt.Errorf("failed to check like: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("user already liked this post")
	}

	// Помечаем, что пользователь лайкнул
	if err := s.redis.Set(ctx, userLikedKey, "1", 0).Err(); err != nil {
		return fmt.Errorf("failed to set like: %w", err)
	}

	// Увеличиваем счетчик лайков
	likesKey := fmt.Sprintf(keyPostLikes, postID)
	likes, err := s.redis.Incr(ctx, likesKey).Result()
	if err != nil {
		return fmt.Errorf("failed to increment likes: %w", err)
	}

	// Публикуем событие лайка
	event := PostLikedEvent{
		PostID: postID,
		UserID: userID,
		Likes:  int(likes),
	}
	if err := s.publishEvent("post.liked", event); err != nil {
		s.logger.Error("failed to publish post.liked event", err, nil)
	}

	return nil
}

// UnlikePost убирает лайк с поста
func (s *MetricsService) UnlikePost(ctx context.Context, userID, postID int) error {
	userLikedKey := fmt.Sprintf(keyUserLikedPost, userID, postID)
	exists, err := s.redis.Exists(ctx, userLikedKey).Result()
	if err != nil {
		return fmt.Errorf("failed to check like: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("user hasn't liked this post")
	}

	// Удаляем отметку лайка
	if err := s.redis.Del(ctx, userLikedKey).Err(); err != nil {
		return fmt.Errorf("failed to remove like: %w", err)
	}

	// Уменьшаем счетчик лайков
	likesKey := fmt.Sprintf(keyPostLikes, postID)
	likes, err := s.redis.Decr(ctx, likesKey).Result()
	if err != nil {
		return fmt.Errorf("failed to decrement likes: %w", err)
	}

	// Публикуем событие удаления лайка
	event := PostUnlikedEvent{
		PostID: postID,
		UserID: userID,
		Likes:  int(likes),
	}
	if err := s.publishEvent("post.unliked", event); err != nil {
		s.logger.Error("failed to publish post.unliked event", err, nil)
	}

	return nil
}

// GetLikes возвращает количество лайков поста
func (s *MetricsService) GetLikes(ctx context.Context, postID int) (int, error) {
	key := fmt.Sprintf(keyPostLikes, postID)
	likes, err := s.redis.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get likes: %w", err)
	}
	return likes, nil
}

// IsLiked проверяет, лайкнул ли пользователь пост
func (s *MetricsService) IsLiked(ctx context.Context, userID, postID int) (bool, error) {
	key := fmt.Sprintf(keyUserLikedPost, userID, postID)
	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check like: %w", err)
	}
	return exists > 0, nil
}

// GetMetrics возвращает метрики поста (лайки и просмотры)
func (s *MetricsService) GetMetrics(ctx context.Context, postID int) (likes int, views int, err error) {
	likesKey := fmt.Sprintf(keyPostLikes, postID)
	viewsKey := fmt.Sprintf(keyPostViews, postID)

	pipe := s.redis.Pipeline()
	likesCmd := pipe.Get(ctx, likesKey)
	viewsCmd := pipe.Get(ctx, viewsKey)

	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return 0, 0, fmt.Errorf("failed to get metrics: %w", err)
	}

	if likesCmd.Err() == nil {
		likes, _ = strconv.Atoi(likesCmd.Val())
	}
	if viewsCmd.Err() == nil {
		views, _ = strconv.Atoi(viewsCmd.Val())
	}

	return likes, views, nil
}

func (s *MetricsService) publishEvent(topic string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	if err := s.publisher.Publish(topic, msg); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}
