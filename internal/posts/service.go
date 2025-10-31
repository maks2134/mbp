package posts

import (
	"context"
	"encoding/json"
	"fmt"
	"mpb/pkg/errors_constant"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type PostsRepositoryInterface interface {
	Save(ctx context.Context, post *Post) error
	FindByID(ctx context.Context, id int) (*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, f PostFilter) ([]Post, error)
}

type PostsService struct {
	repo      PostsRepositoryInterface
	publisher message.Publisher
	logger    watermill.LoggerAdapter
}

func NewPostsService(repo *PostsRepository, publisher message.Publisher, logger watermill.LoggerAdapter) *PostsService {
	return &PostsService{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}

func (s *PostsService) CreatePost(ctx context.Context, userID int, title, description, tag string) (*Post, error) {
	title = strings.TrimSpace(title)
	if len(title) < 3 {
		return nil, errors_constant.InvalidTitle
	}

	post := &Post{
		UserID:       userID,
		Title:        title,
		Description:  description,
		Tag:          tag,
		Like:         0,
		CountViewers: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Save(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	event := PostCreatedEvent{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		CreatedAt: post.CreatedAt,
	}

	payload, _ := json.Marshal(event)
	msg := message.NewMessage(watermill.NewUUID(), payload)

	if err := s.publisher.Publish("post.created", msg); err != nil {
		s.logger.Error("failed to publish post.created event", err, nil)
	}

	return post, nil
}

func (s *PostsService) GetPostByID(ctx context.Context, id int) (*Post, error) {
	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors_constant.PostNotFound
	}
	return post, nil
}

func (s *PostsService) UpdatePost(ctx context.Context, userID, postID int, title, description, tag string) (*Post, error) {
	post, err := s.repo.FindByID(ctx, postID)
	if err != nil {
		return nil, errors_constant.PostNotFound
	}

	if post.UserID != userID {
		return nil, errors_constant.UserNotAuthorized
	}

	post.Title = strings.TrimSpace(title)
	post.Description = description
	post.Tag = tag
	post.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	// TODO: publish PostUpdated event
	return post, nil
}

func (s *PostsService) DeletePost(ctx context.Context, userID, postID int) error {
	post, err := s.repo.FindByID(ctx, postID)
	if err != nil {
		return errors_constant.PostNotFound
	}

	if post.UserID != userID {
		return errors_constant.UserNotFound
	}

	if err := s.repo.Delete(ctx, postID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	// TODO: publish PostDeleted event
	return nil
}

func (s *PostsService) ListPosts(ctx context.Context, f PostFilter) ([]Post, error) {
	posts, err := s.repo.List(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	return posts, nil
}
