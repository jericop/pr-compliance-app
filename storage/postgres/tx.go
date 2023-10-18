package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseTransactionFactory interface {
	ExecWithTx(context.Context, func(Querier) error) error
}

type postgresTransactionFactory struct {
	pool *pgxpool.Pool
}

func NewPostgresTransactionFactory(p *pgxpool.Pool) *postgresTransactionFactory {
	return &postgresTransactionFactory{pool: p}
}

func (p *postgresTransactionFactory) ExecWithTx(ctx context.Context, fn func(Querier) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("tx begin err: %v", err)
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	if cErr := tx.Commit(ctx); cErr != nil {
		return fmt.Errorf("tx commit err: %v, fn err: %v", cErr, err)
	}

	return err
}
