package gormtx

import (
	"context"

	"gorm.io/gorm"
)

type txCtxKey struct{}

func setTxCtx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txCtxKey{}, tx)
}

func getTxCtx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txCtxKey{}).(*gorm.DB)
	return tx, ok
}
