package dtos

import (
	"github.com/Govorov1705/ozon-test/internal/models"
	"github.com/google/uuid"
)

type CommentWithReplies struct {
	models.Comment
	Replies []*CommentWithReplies
}

type CreateCommentRequest struct {
	PostID  uuid.UUID `validate:"required"`
	UserID  uuid.UUID `validate:"required"`
	ReplyTo *uuid.UUID
	Content string `validate:"required,max=2000"`
}
