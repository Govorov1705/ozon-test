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

type CommentsRepository struct {
	*BaseRepository
}

func NewCommentsRepository(pool *pgxpool.Pool) repositories.CommentsRepository {
	return &CommentsRepository{
		BaseRepository: NewBaseRepository(pool),
	}
}

func (r *CommentsRepository) Add(ctx context.Context, postID, userID uuid.UUID, rootID, replyTo *uuid.UUID, content string) (*models.Comment, error) {
	comment := models.Comment{}

	stmt := `
		INSERT INTO comments(id, post_id, user_id, root_id, reply_to, content) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, post_id, user_id, root_id, reply_to, content, created_at;
	`

	commentID := uuid.New()
	if rootID == nil {
		rootID = &commentID
	}

	querier := r.GetQuerier(ctx)
	row := querier.QueryRow(ctx, stmt, commentID, postID, userID, rootID, replyTo, content)

	err := row.Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.RootID,
		&comment.ReplyTo,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		logger.Logger.Error("error scanning row", zap.Error(err))
		return nil, errs.ErrInternal
	}

	return &comment, nil
}

func (r *CommentsRepository) GetByID(ctx context.Context, commentID uuid.UUID, forUpdate bool) (*models.Comment, error) {
	comment := models.Comment{}

	query := `
		SELECT id, post_id, user_id, root_id, reply_to, content, created_at 
		FROM comments
		WHERE id = $1`
	if forUpdate {
		query += " FOR UPDATE"
	}
	query += ";"

	querier := r.GetQuerier(ctx)
	row := querier.QueryRow(ctx, query, commentID)

	err := row.Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.RootID,
		&comment.ReplyTo,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("comment %w", errs.ErrNotFound)
		}
		logger.Logger.Error("error scanning row", zap.Error(err))
		return nil, errs.ErrInternal
	}

	return &comment, nil
}

func (r *CommentsRepository) GetRootCommentsByPostID(ctx context.Context, postID uuid.UUID, limit, offset *int32) ([]*models.Comment, error) {
	comments := []*models.Comment{}

	l := int32(10)
	o := int32(0)

	if limit != nil {
		l = *limit
	}
	if offset != nil {
		o = *offset
	}

	query := `
		SELECT id, post_id, user_id, root_id, reply_to, content, created_at
		FROM comments
		WHERE post_id = $1 AND reply_to is NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	querier := r.GetQuerier(ctx)
	rows, err := querier.Query(ctx, query, postID, l, o)
	if err != nil {
		logger.Logger.Error("error during query", zap.Error(err))
		return nil, errs.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		comment := models.Comment{}

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.RootID,
			&comment.ReplyTo,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			logger.Logger.Error("error scanning row", zap.Error(err))
			return nil, errs.ErrInternal
		}

		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		logger.Logger.Error("error during row iteration", zap.Error(err))
		return nil, errs.ErrInternal
	}

	return comments, nil
}

func (r *CommentsRepository) GetChildrenCommentsByRootIDs(ctx context.Context, rootIDs []*uuid.UUID) ([]*models.Comment, error) {
	comments := []*models.Comment{}

	query := `
		SELECT id, post_id, user_id, root_id, reply_to, content, created_at
		FROM comments
		WHERE root_id = ANY($1)
		ORDER BY created_at DESC;
	`

	querier := r.GetQuerier(ctx)
	rows, err := querier.Query(ctx, query, rootIDs)
	if err != nil {
		logger.Logger.Error("error during query", zap.Error(err))
		return nil, errs.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		comment := models.Comment{}

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.RootID,
			&comment.ReplyTo,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			logger.Logger.Error("error scanning row", zap.Error(err))
			return nil, errs.ErrInternal
		}

		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		logger.Logger.Error("error during row iteration", zap.Error(err))
		return nil, errs.ErrInternal
	}

	return comments, nil
}
