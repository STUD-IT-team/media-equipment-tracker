package gormtx

import (
	"context"
	"errors"

	"media-equipment-tracker/pkg/txmanager"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type GormTxManager struct {
	db         *gorm.DB
	retryCount int
}

func NewGormTxManager(db *gorm.DB, retryCount int) *GormTxManager {
	return &GormTxManager{
		db:         db,
		retryCount: retryCount,
	}
}

var _ txmanager.TxManager = (*GormTxManager)(nil)

// WithinTx выполняет fn в рамках транзакции.
//
// Если транзакция уже присутствует в ctx, функция выполняется
// без создания новой транзакции.
//
// Для ошибок сериализации и дедлоков выполняется повторная попытка
// выполнения транзакции.
func (m *GormTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := getTxCtx(ctx); ok {
		return fn(ctx)
	}
	var lastErr error

	for i := 0; i <= m.retryCount; i++ {
		err := m.runTx(ctx, fn)
		if err == nil {
			return nil
		}

		// Gorm используей pgx под капотом
		var postgresErr *pgconn.PgError
		if !errors.Is(err, postgresErr) {
			return err
		}
		switch postgresErr.Code {
		case "40001", "40P01": // serialization_failure или deadlock_detected
			lastErr = err
			continue
		default:
			return err
		}
	}
	return txmanager.WrapRetryExceeded(lastErr)
}

func (m *GormTxManager) runTx(ctx context.Context, fn func(ctx context.Context) error) error {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		ctx = setTxCtx(ctx, tx)

		err := fn(ctx)
		if err != nil {
			return txmanager.WrapTransactionClosureError(err)
		}
		return nil
	})
	if err != nil {
		return txmanager.WrapCommitError(err)
	}

	return nil
}
