package pgtx

import (
	"context"
	"errors"
	"media-equipment-tracker/pkg/txmanager"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgTxManager реализует txmanager.TxManager для PostgreSQL через pgx.
type PgTxManager struct {
	pool       *pgxpool.Pool
	retryCount int
}

func newTxManager(pool *pgxpool.Pool, retryCount int) *PgTxManager {
	return &PgTxManager{
		pool:       pool,
		retryCount: retryCount,
	}
}

var _ txmanager.TxManager = (*PgTxManager)(nil)

// WithinTx выполняет fn в рамках транзакции.
//
// Если транзакция уже присутствует в ctx, функция выполняется
// без создания новой транзакции.
//
// Для ошибок сериализации и дедлоков выполняется повторная попытка
// выполнения транзакции.
func (m *PgTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := getTxCtx(ctx); ok {
		return fn(ctx)
	}
	var lastErr error

	for i := 0; i <= m.retryCount; i++ {
		err := m.runTx(ctx, fn)
		if err == nil {
			return nil
		}

		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return err
		}
		switch pgErr.Code {
		case "40001", "40P01": // serialization_failure или deadlock_detected
			lastErr = err
			continue
		default:
			return err
		}
	}
	return txmanager.WrapRetryExceeded(lastErr)
}

func (m *PgTxManager) runTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return txmanager.WrapBeginError(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	ctx = setTxCtx(ctx, tx)

	err = fn(ctx)
	if err != nil {
		return txmanager.WrapTransactionClosureError(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return txmanager.WrapCommitError(err)
	}

	return nil
}
