package handler

import (
	"context"
	"mpb/internal/comments"
	common_proto "mpb/proto/common"
	posts_proto "mpb/proto/posts"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommentsHandler struct {
	posts_proto.UnimplementedCommentsServiceServer
	service *comments.CommentsService
}

func NewCommentsHandler(service *comments.CommentsService) *CommentsHandler {
	return &CommentsHandler{
		service: service,
	}
}

func (h *CommentsHandler) CreateComment(ctx context.Context, req *posts_proto.CreateCommentRequest) (*posts_proto.CommentResponse, error) {
	comment, err := h.service.CreateComment(ctx, int(req.PostId), int(req.UserId), req.Content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create comment: %v", err)
	}

	return &posts_proto.CommentResponse{
		Comment: toProtoComment(comment),
	}, nil
}

func (h *CommentsHandler) GetComment(ctx context.Context, req *posts_proto.GetCommentRequest) (*posts_proto.CommentResponse, error) {
	comment, err := h.service.GetCommentByID(ctx, int(req.CommentId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "comment not found: %v", err)
	}

	return &posts_proto.CommentResponse{
		Comment: toProtoComment(comment),
	}, nil
}

func (h *CommentsHandler) ListComments(ctx context.Context, req *posts_proto.ListCommentsRequest) (*posts_proto.ListCommentsResponse, error) {
	commentList, err := h.service.ListComments(ctx, int(req.PostId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list comments: %v", err)
	}

	protoComments := make([]*posts_proto.Comment, len(commentList))
	for i, c := range commentList {
		protoComments[i] = toProtoComment(&c)
	}

	return &posts_proto.ListCommentsResponse{
		Comments: protoComments,
		Total:    int32(len(protoComments)),
	}, nil
}

func (h *CommentsHandler) UpdateComment(ctx context.Context, req *posts_proto.UpdateCommentRequest) (*posts_proto.CommentResponse, error) {
	comment, err := h.service.UpdateComment(ctx, int(req.UserId), int(req.CommentId), req.Content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update comment: %v", err)
	}

	return &posts_proto.CommentResponse{
		Comment: toProtoComment(comment),
	}, nil
}

func (h *CommentsHandler) DeleteComment(ctx context.Context, req *posts_proto.DeleteCommentRequest) (*common_proto.Empty, error) {
	err := h.service.DeleteComment(ctx, int(req.UserId), int(req.CommentId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete comment: %v", err)
	}

	return &common_proto.Empty{}, nil
}

func toProtoComment(c *comments.Comment) *posts_proto.Comment {
	return &posts_proto.Comment{
		Id:        int32(c.ID),
		PostId:    int32(c.PostID),
		UserId:    int32(c.UserID),
		Content:   c.Text,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
