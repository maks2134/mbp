package users

import (
	"context"
	"mpb/internal/posts"
	"mpb/internal/user"
	"mpb/pkg/errors_constant"
)

type UsersService struct {
	repo      *UsersRepository
	postsRepo posts.PostsRepositoryInterface
}

func NewUsersService(repo *UsersRepository, postsRepo posts.PostsRepositoryInterface) *UsersService {
	return &UsersService{
		repo:      repo,
		postsRepo: postsRepo,
	}
}

func (s *UsersService) GetUserProfile(ctx context.Context, userID int) (*UserProfile, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors_constant.UserNotFound
	}

	postsCount, err := s.repo.GetPostsCount(ctx, userID)
	if err != nil {
		postsCount = 0
	}

	attachmentsCount, err := s.repo.GetAttachmentsCount(ctx, userID)
	if err != nil {
		attachmentsCount = 0
	}

	return &UserProfile{
		User:             *u,
		PostsCount:       postsCount,
		AttachmentsCount: attachmentsCount,
	}, nil
}

func (s *UsersService) GetUserPosts(ctx context.Context, userID int, filter posts.PostFilter) ([]posts.Post, error) {
	_, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors_constant.UserNotFound
	}

	filter.UserID = &userID
	filter.OnlyActive = true

	return s.postsRepo.List(ctx, filter)
}

func (s *UsersService) ListUsers(ctx context.Context, filter UserFilter) ([]user.User, error) {
	return s.repo.List(ctx, filter)
}
