package repositories

import (
	"context"

	"github.com/Govorov1705/ozon-test/internal/models"
	"github.com/google/uuid"
)

type CommentsRepository interface {
	Add(ctx context.Context, postID, userID uuid.UUID, rootID, replyTo *uuid.UUID, content string) (*models.Comment, error)
	GetByID(ctx context.Context, commentID uuid.UUID, forUpdate bool) (*models.Comment, error)
	GetRootCommentsByPostID(ctx context.Context, postID uuid.UUID, limit, offset *int32) ([]*models.Comment, error)
	GetChildrenCommentsByRootIDs(ctx context.Context, rootIDs []*uuid.UUID) ([]*models.Comment, error)
}
