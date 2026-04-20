package gormtx

import (
	"context"

	"gorm.io/gorm"
)

// DBGetter -- обёртка над gorm.DB, при использовании с GormTxManager
// позволяет получать общую транзакцию из контекста или создавать новую.
// Используется на уровне доступа к данным, для унификации работы в транзакции и без неё.

type DBGetter struct {
	db *gorm.DB
}

func newDBGetter(db *gorm.DB) *DBGetter {
	return &DBGetter{db: db}
}

func (g *DBGetter) GetDB(ctx context.Context) (*gorm.DB, error) {
	tx, ok := getTxCtx(ctx)
	if ok {
		return tx, nil
	}
	return g.db, nil
}
