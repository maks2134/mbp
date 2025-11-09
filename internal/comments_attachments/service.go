package comments_attachments

import "context"

type CommentAttachmentsService struct {
	repo *CommentAttachmentsRepository
}

func NewCommentAttachmentsService(repo *CommentAttachmentsRepository) *CommentAttachmentsService {
	return &CommentAttachmentsService{repo: repo}
}

func (s *CommentAttachmentsService) CreateAttachment(ctx context.Context, att *CommentAttachment) error {
	return s.repo.Create(ctx, att)
}

func (s *CommentAttachmentsService) ListAttachments(ctx context.Context, commentID int) ([]CommentAttachment, error) {
	return s.repo.ListByComment(ctx, commentID)
}

func (s *CommentAttachmentsService) DeleteAttachment(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
