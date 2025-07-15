package repositories

import (
	"context"

	"github.com/Govorov1705/ozon-test/internal/models"
)

type UsersRepository interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Add(ctx context.Context, username, hashedPassword string) (*models.User, error)
}
