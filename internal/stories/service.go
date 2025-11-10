package stories

import (
	"context"
	"fmt"
	"mpb/pkg/errors_constant"
	"time"
)

type StoriesService struct {
	repo *StoriesRepository
}

func NewStoriesService(repo *StoriesRepository) *StoriesService {
	return &StoriesService{repo: repo}
}

func (s *StoriesService) CreateStory(ctx context.Context, userID int, fileURL, fileType string) (*Story, error) {
	story := &Story{
		UserID:     userID,
		FileURL:    fileURL,
		FileType:   fileType,
		ViewsCount: 0,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.Create(ctx, story); err != nil {
		return nil, fmt.Errorf("failed to create story: %w", err)
	}

	return story, nil
}

func (s *StoriesService) GetStory(ctx context.Context, storyID int, viewerUserID *int) (*Story, bool, error) {
	story, err := s.repo.FindByID(ctx, storyID)
	if err != nil {
		return nil, false, errors_constant.UserNotFound // Можно создать отдельную ошибку StoryNotFound
	}

	var isViewed bool
	if viewerUserID != nil {
		isViewed, _ = s.repo.HasUserViewed(ctx, storyID, *viewerUserID)
	}

	return story, isViewed, nil
}

func (s *StoriesService) ViewStory(ctx context.Context, storyID, userID int) error {
	story, err := s.repo.FindByID(ctx, storyID)
	if err != nil {
		return errors_constant.UserNotFound
	}

	if story.UserID == userID {
		return nil
	}

	hasViewed, err := s.repo.HasUserViewed(ctx, storyID, userID)
	if err != nil {
		return fmt.Errorf("failed to check view: %w", err)
	}

	if !hasViewed {
		if err := s.repo.IncrementViews(ctx, storyID); err != nil {
			return fmt.Errorf("failed to increment views: %w", err)
		}
		if err := s.repo.RecordView(ctx, storyID, userID); err != nil {
			return fmt.Errorf("failed to record view: %w", err)
		}
	}

	return nil
}

func (s *StoriesService) ListUserStories(ctx context.Context, userID int) ([]Story, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *StoriesService) ListActiveStories(ctx context.Context, excludeUserID *int) ([]Story, error) {
	return s.repo.ListActive(ctx, excludeUserID)
}

func (s *StoriesService) DeleteStory(ctx context.Context, storyID, userID int) error {
	return s.repo.Delete(ctx, storyID, userID)
}

func (s *StoriesService) CleanupExpired(ctx context.Context) error {
	return s.repo.DeleteExpired(ctx)
}
