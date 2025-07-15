package repositories

import (
	"context"

	"github.com/Govorov1705/ozon-test/internal/storages/postgresql"
	"github.com/Govorov1705/ozon-test/internal/transactions"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseRepository struct {
	Pool *pgxpool.Pool
}

func NewBaseRepository(pool *pgxpool.Pool) *BaseRepository {
	return &BaseRepository{Pool: pool}
}

func (r *BaseRepository) GetQuerier(ctx context.Context) postgresql.Querier {
	if tx, ok := transactions.GetTxFromContext(ctx); ok {
		if pgxTx, ok := tx.(pgx.Tx); ok {
			return pgxTx
		}
	}
	return r.Pool
}
