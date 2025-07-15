package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/Govorov1705/ozon-test/internal/errs"
	"github.com/Govorov1705/ozon-test/internal/logger"
	"github.com/Govorov1705/ozon-test/internal/models"
	"github.com/Govorov1705/ozon-test/internal/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostsRepository struct {
	*BaseRepository
}

func NewPostsRepository(pool *pgxpool.Pool) repositories.PostsRepository {
	return &PostsRepository{
		BaseRepository: NewBaseRepository(pool),
	}
}

func (r *PostsRepository) Add(ctx context.Context, userID uuid.UUID, title, content string, areCommentsAllowed bool) (*models.Post, error) {
	post := models.Post{}

	stmt := `
		INSERT INTO posts(user_id, title, content, are_comments_allowed) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, content, are_comments_allowed, created_at;
	`

	querier := r.GetQuerier(ctx)
	row := querier.QueryRow(ctx, stmt, userID, title, content, areCommentsAllowed)

	err := row.Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.AreCommentsAllowed,
		&post.CreatedAt,
	)
	if err != nil {
		logger.Logger.Error("error scanning row", zap.Error(err))
		return nil, errs.ErrInternal
	}

	return &post, nil
}

func (r *PostsRepository) GetByID(ctx context.Context, postID uuid.UUID, forUpdate bool) (*models.Post, error) {
	post := models.Post{}

	query := `
		SELECT id, user_id, title, content, are_comments_allowed, created_at 
		FROM posts
		WHERE id = $1`
	if forUpdate {
		query += " FOR UPDATE"
	}
	query += ";"

	querier := r.GetQuerier(ctx)
	row := querier.QueryRow(ctx, query, postID)

	err := row.Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.AreCommentsAllowed,
		&post.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("post %w", errs.ErrNotFound)
		}
		logger.Logger.Error("error scanning row", zap.Error(err))
		return nil, errs.ErrInternal
	}

	return &post, nil
}

func (r *PostsRepository) GetAll(ctx context.Context) ([]*models.Post, error) {
	posts := []*models.Post{}

	query := `
		SELECT id, user_id, title, content, are_comments_allowed, created_at
		FROM posts
		ORDER BY created_at DESC;
	`

	querier := r.GetQuerier(ctx)
	rows, err := querier.Query(ctx, query)
	if err != nil {
		logger.Logger.Error("error during query", zap.Error(err))
		return nil, errs.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		post := models.Post{}

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.AreCommentsAllowed,
			&post.CreatedAt,
		)
		if err != nil {
			logger.Logger.Error("error scanning row", zap.Error(err))
			return nil, errs.ErrInternal
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		logger.Logger.Error("error during row iteration", zap.Error(err))
		return nil, errs.ErrInternal
	}

	return posts, nil
}

func (r *PostsRepository) DisableComments(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	post := models.Post{}

	stmt := `
		UPDATE posts
		SET are_comments_allowed = false
		WHERE id = $1
		RETURNING id, user_id, title, content, are_comments_allowed, created_at;
	`

	querier := r.GetQuerier(ctx)
	row := querier.QueryRow(ctx, stmt, postID)

	err := row.Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.AreCommentsAllowed,
		&post.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("post %w", errs.ErrNotFound)
		}
		logger.Logger.Error("error scanning row", zap.Error(err))
		return nil, errs.ErrInternal
	}

	return &post, nil
}

func (r *PostsRepository) EnableComments(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	post := models.Post{}

	stmt := `
		UPDATE posts
		SET are_comments_allowed = true
		WHERE id = $1
		RETURNING id, user_id, title, content, are_comments_allowed, created_at;
	`

	querier := r.GetQuerier(ctx)
	row := querier.QueryRow(ctx, stmt, postID)

	err := row.Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.AreCommentsAllowed,
		&post.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("post %w", errs.ErrNotFound)
		}
		logger.Logger.Error("error scanning row", zap.Error(err))
		return nil, errs.ErrInternal
	}

	return &post, nil
}
