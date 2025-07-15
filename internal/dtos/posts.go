package dtos

import (
	"github.com/Govorov1705/ozon-test/internal/models"
	"github.com/google/uuid"
)

type PostWithComments struct {
	Post     *models.Post
	Comments []*CommentWithReplies
}

type CreatePostRequest struct {
	UserID             uuid.UUID `validate:"required"`
	Title              string    `validate:"required,max=100"`
	Content            string    `validate:"required,max=2000"`
	AreCommentsAllowed *bool
}

type GetPostWithCommentsRequest struct {
	PostID uuid.UUID `validate:"required"`
	Limit  *int32    `validate:"gt=0"`
	Offset *int32    `validate:"gte=0"`
}
