package postgresql

import (
	"context"

	"github.com/Govorov1705/ozon-test/internal/transactions"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxpoolTxStarter struct {
	pool *pgxpool.Pool
}

func NewPgxpoolTxStarter(pool *pgxpool.Pool) *PgxpoolTxStarter {
	return &PgxpoolTxStarter{pool: pool}
}

func (p *PgxpoolTxStarter) Begin(ctx context.Context) (transactions.Tx, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
