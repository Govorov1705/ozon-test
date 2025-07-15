package repositories

import (
	"context"

	"github.com/Govorov1705/ozon-test/internal/models"
	"github.com/google/uuid"
)

type PostsRepository interface {
	Add(ctx context.Context, userID uuid.UUID, title, content string, areCommentsAllowed bool) (*models.Post, error)
	GetByID(ctx context.Context, postID uuid.UUID, forUpdate bool) (*models.Post, error)
	GetAll(ctx context.Context) ([]*models.Post, error)
	DisableComments(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	EnableComments(ctx context.Context, postID uuid.UUID) (*models.Post, error)
}
