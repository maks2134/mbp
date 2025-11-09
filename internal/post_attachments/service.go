package post_attachments

import "context"

type PostAttachmentsService struct {
	repo *PostAttacmentsRepository
}

func NewPostAttachmentsService(repo *PostAttacmentsRepository) *PostAttachmentsService {
	return &PostAttachmentsService{repo: repo}
}

func (s *PostAttachmentsService) CreateAttachment(ctx context.Context, att *PostAttachment) error {
	return s.repo.Create(ctx, att)
}

func (s *PostAttachmentsService) ListAttachments(ctx context.Context, postID int) ([]PostAttachment, error) {
	return s.repo.ListByPost(ctx, postID)
}

func (s *PostAttachmentsService) DeleteAttachment(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
