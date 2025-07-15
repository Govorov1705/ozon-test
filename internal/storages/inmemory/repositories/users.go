package repositories

import (
	"context"
	"sync"

	"github.com/Govorov1705/ozon-test/internal/errs"
	"github.com/Govorov1705/ozon-test/internal/models"
	"github.com/Govorov1705/ozon-test/internal/repositories"
	"github.com/google/uuid"
)

type InMemoryUsersRepository struct {
	mu    sync.RWMutex
	users map[string]*models.User
}

func NewUsersRepository() repositories.UsersRepository {
	return &InMemoryUsersRepository{
		users: make(map[string]*models.User),
	}
}

func (r *InMemoryUsersRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[username]
	if !ok {
		return nil, errs.ErrNotFound
	}

	return user, nil
}

func (r *InMemoryUsersRepository) Add(ctx context.Context, username, hashedPassword string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.users[username]
	if ok {
		return nil, errs.ErrAlreadyExists
	}

	user := &models.User{
		ID:             uuid.New(),
		Username:       username,
		HashedPassword: hashedPassword,
	}

	r.users[username] = user

	return user, nil
}
