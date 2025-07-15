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

type InMemoryCommentsRepository struct {
	mu       sync.RWMutex
	comments map[uuid.UUID]*models.Comment
}

func NewCommentsRepository() repositories.CommentsRepository {
	return &InMemoryCommentsRepository{
		comments: make(map[uuid.UUID]*models.Comment),
	}
}

func (r *InMemoryCommentsRepository) Add(ctx context.Context, postID, userID uuid.UUID, rootID, replyTo *uuid.UUID, content string) (*models.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	commentID := uuid.New()
	if rootID == nil {
		rootID = &commentID
	}

	comment := &models.Comment{
		ID:        commentID,
		PostID:    postID,
		UserID:    userID,
		RootID:    *rootID,
		ReplyTo:   replyTo,
		Content:   content,
		CreatedAt: time.Now(),
	}

	r.comments[comment.ID] = comment

	return comment, nil
}

func (r *InMemoryCommentsRepository) GetByID(ctx context.Context, commentID uuid.UUID, forUpdate bool) (*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, ok := r.comments[commentID]
	if !ok {
		return nil, errs.ErrNotFound
	}

	return comment, nil
}

func (r *InMemoryCommentsRepository) GetRootCommentsByPostID(ctx context.Context, postID uuid.UUID, limit, offset *int32) ([]*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var rootsComments []*models.Comment

	for _, comment := range r.comments {
		if comment.PostID == postID && comment.ReplyTo == nil {
			rootsComments = append(rootsComments, comment)
		}
	}

	sort.Slice(rootsComments, func(i, j int) bool {
		return rootsComments[i].CreatedAt.After(rootsComments[j].CreatedAt)
	})

	start := int32(0)
	if offset != nil {
		start = *offset
	}

	end := int32(len(rootsComments))
	if limit != nil && start+*limit < end {
		end = start + *limit
	}

	if start > end {
		return []*models.Comment{}, nil
	}

	return rootsComments[start:end], nil
}

func (r *InMemoryCommentsRepository) GetChildrenCommentsByRootIDs(ctx context.Context, rootIDs []*uuid.UUID) ([]*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comments := []*models.Comment{}
	IDSet := map[uuid.UUID]struct{}{}

	for _, ID := range rootIDs {
		IDSet[*ID] = struct{}{}
	}

	for _, comment := range r.comments {
		_, ok := IDSet[comment.RootID]
		if ok {
			comments = append(comments, comment)
		}
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.After(comments[j].CreatedAt)
	})

	return comments, nil
}
