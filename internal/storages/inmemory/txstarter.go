package inmemory

import (
	"context"

	"github.com/Govorov1705/ozon-test/internal/transactions"
)

type InMemoryTx struct{}

func (i *InMemoryTx) Commit(ctx context.Context) error {
	return nil
}

func (i *InMemoryTx) Rollback(ctx context.Context) error {
	return nil
}

type InMemoryTxStarter struct{}

func (i *InMemoryTxStarter) Begin(ctx context.Context) (transactions.Tx, error) {
	return &InMemoryTx{}, nil
}
