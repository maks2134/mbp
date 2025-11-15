package comments

import (
	"context"
	"fmt"
	"mpb/pkg/errors_constant"
	"time"
)

type CommentsRepositoryInterface interface {
	Create(ctx context.Context, c *Comment) error
	Update(ctx context.Context, c *Comment) error
	Delete(ctx context.Context, commentID int) error
	List(ctx context.Context, postID int) ([]Comment, error)
	FindCommentByID(ctx context.Context, commentID int) (*Comment, error)
}

type CommentsService struct {
	repo CommentsRepositoryInterface
}

func NewCommentsService(repo CommentsRepositoryInterface) *CommentsService {
	return &CommentsService{repo: repo}
}

func (s *CommentsService) CreateComment(ctx context.Context, postID, userID int, text string) (*Comment, error) {
	if len(text) == 0 {
		return nil, errors_constant.InvalidCommentText
	}

	comment := &Comment{
		PostID:    postID,
		UserID:    userID,
		Text:      text,
		Like:      0,
		Blocked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

func (s *CommentsService) UpdateComment(ctx context.Context, userID, commentID int, newText string) (*Comment, error) {
	comment, err := s.repo.FindCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find comment: %w", err)
	}

	if comment.DeletedAt != nil {
		return nil, errors_constant.CommentDeleted
	}

	if comment.UserID != userID {
		return nil, errors_constant.UserNotAuthorized
	}

	if len(newText) == 0 {
		return nil, errors_constant.InvalidCommentText
	}

	comment.Text = newText
	comment.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return comment, nil
}

func (s *CommentsService) DeleteComment(ctx context.Context, userID, commentID int) error {
	comment, err := s.repo.FindCommentByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to find comment: %w", err)
	}

	if comment.UserID != userID {
		return errors_constant.UserNotAuthorized
	}

	if err := s.repo.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func (s *CommentsService) ListComments(ctx context.Context, postID int) ([]Comment, error) {
	comments, err := s.repo.List(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}

	return comments, nil
}

func (s *CommentsService) GetCommentByID(ctx context.Context, commentID int) (*Comment, error) {
	comment, err := s.repo.FindCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find comment: %w", err)
	}

	if comment.DeletedAt != nil {
		return nil, errors_constant.CommentDeleted
	}

	return comment, nil
}

func (s *CommentsService) LikeComment(ctx context.Context, userID, commentID int) (*Comment, error) {
	comment, err := s.repo.FindCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find comment: %w", err)
	}

	if comment.DeletedAt != nil {
		return nil, errors_constant.CommentDeleted
	}

	comment.Like += 1
	comment.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to like comment: %w", err)
	}

	return comment, nil
}

func (s *CommentsService) BlockComment(ctx context.Context, commentID int, block bool) (*Comment, error) {
	comment, err := s.repo.FindCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find comment: %w", err)
	}

	comment.Blocked = block
	comment.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to update blocked state: %w", err)
	}

	return comment, nil
}
