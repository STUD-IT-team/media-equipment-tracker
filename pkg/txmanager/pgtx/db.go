package pgtx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB -- обёртка на pgxpool, при использовании с PgTxManager
// позволяет получать общую транзакцию из контекста или создавать новую.
// Используется на уровне доступа к данным, для унификации работы в транзакции и без неё.
type DB struct {
	pool *pgxpool.Pool
}

func newDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool: pool,
	}
}

// Querier описывает минимальный интерфейс для выполнения SQL-запросов с помощью pgx.
type Querier interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

// GetConn возвращает соединение для выполнения запросов.
//
// Если в ctx присутствует транзакция — возвращается она (т.е. код выполняется WithinTx),
// иначе используется пул соединений.
func (db *DB) GetConn(ctx context.Context) (Querier, error) {
	tx, ok := getTxCtx(ctx)
	if ok {
		return tx, nil
	}

	return db.pool, nil
}
