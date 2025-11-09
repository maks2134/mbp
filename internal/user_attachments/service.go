package user_attachments

import (
	"context"
	"mpb/pkg/errors_constant"
)

type UserAttachmentsService struct {
	repo *UserAttachmentsRepository
}

func NewUserAttachmentsService(repo *UserAttachmentsRepository) *UserAttachmentsService {
	return &UserAttachmentsService{repo: repo}
}

func (s *UserAttachmentsService) CreateAttachment(ctx context.Context, att *UserAttachment) error {
	return s.repo.Create(ctx, att)
}

func (s *UserAttachmentsService) ListAttachments(ctx context.Context, userID int) ([]UserAttachment, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *UserAttachmentsService) DeleteAttachment(ctx context.Context, id int, userID int) error {
	att, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors_constant.UserNotFound
	}

	if att.UserID != userID {
		return errors_constant.UserNotAuthorized
	}

	return s.repo.Delete(ctx, id)
}
