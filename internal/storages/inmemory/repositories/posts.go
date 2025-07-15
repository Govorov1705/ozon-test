package repositories

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/Govorov1705/ozon-test/internal/errs"
	"github.com/Govorov1705/ozon-test/internal/models"
	"github.com/Govorov1705/ozon-test/internal/repositories"
	"github.com/google/uuid"
)

type InMemoryPostsRepository struct {
	mu    sync.RWMutex
	posts map[uuid.UUID]*models.Post
}

func NewPostsRepository() repositories.PostsRepository {
	return &InMemoryPostsRepository{
		posts: make(map[uuid.UUID]*models.Post),
	}
}

func (r *InMemoryPostsRepository) Add(ctx context.Context, userID uuid.UUID, title, content string, areCommentsAllowed bool) (*models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	post := &models.Post{
		ID:                 uuid.New(),
		UserID:             userID,
		Title:              title,
		Content:            content,
		AreCommentsAllowed: areCommentsAllowed,
		CreatedAt:          time.Now(),
	}

	r.posts[post.ID] = post

	return post, nil
}

func (r *InMemoryPostsRepository) GetByID(ctx context.Context, postID uuid.UUID, forUpdate bool) (*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, ok := r.posts[postID]
	if !ok {
		return nil, errs.ErrNotFound
	}

	return post, nil
}

func (r *InMemoryPostsRepository) GetAll(ctx context.Context) ([]*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]*models.Post, 0, len(r.posts))

	for _, post := range r.posts {
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	return posts, nil
}

func (r *InMemoryPostsRepository) DisableComments(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	post, ok := r.posts[postID]
	if !ok {
		return nil, errs.ErrNotFound
	}

	post.AreCommentsAllowed = false

	return post, nil
}

func (r *InMemoryPostsRepository) EnableComments(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	post, ok := r.posts[postID]
	if !ok {
		return nil, errs.ErrNotFound
	}

	post.AreCommentsAllowed = true

	return post, nil
}
