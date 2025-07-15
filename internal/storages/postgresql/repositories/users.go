package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/Govorov1705/ozon-test/internal/errs"
	"github.com/Govorov1705/ozon-test/internal/logger"
	"github.com/Govorov1705/ozon-test/internal/models"
	"github.com/Govorov1705/ozon-test/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UsersRepository struct {
	*BaseRepository
}

func NewUsersRepository(pool *pgxpool.Pool) repositories.UsersRepository {
	return &UsersRepository{
		BaseRepository: NewBaseRepository(pool),
	}
}

func (r *UsersRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := models.User{}

	query := `
		SELECT id, username, hashed_password
		FROM users
		WHERE username = $1;
	`

	querier := r.GetQuerier(ctx)
	row := querier.QueryRow(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user %w", errs.ErrNotFound)
		}
		logger.Logger.Error("error scanning row", zap.Error(err))
		return nil, errs.ErrInternal
	}

	return &user, nil
}

func (r *UsersRepository) Add(ctx context.Context, username, hashedPassword string) (*models.User, error) {
	user := models.User{}

	stmt := `
		INSERT INTO users(username, hashed_password) 
		VALUES ($1, $2)
		RETURNING id, username, hashed_password;
	`

	querier := r.GetQuerier(ctx)
	row := querier.QueryRow(ctx, stmt, username, hashedPassword)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
	)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errs.ErrAlreadyExists
		}
		logger.Logger.Error("error scanning row", zap.Error(err))
		return nil, errs.ErrInternal
	}

	return &user, nil
}
