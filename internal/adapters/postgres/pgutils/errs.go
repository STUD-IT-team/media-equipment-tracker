package pgutils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsUniqueViolationError(err error) bool {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		return pgerr.Code == "23505"
	}
	return false
}
