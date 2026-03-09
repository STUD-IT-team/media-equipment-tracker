package pgtx

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type txCtxKey struct{}

func setTxCtx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey{}, tx)
}

func getTxCtx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txCtxKey{}).(pgx.Tx)
	return tx, ok
}
