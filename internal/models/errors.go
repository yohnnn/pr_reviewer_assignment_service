package models

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInternal      = errors.New("internal error")

	ErrPRMerged     = errors.New("cannot edit merged PR")
	ErrNoCandidates = errors.New("no candidates available")
	ErrNotAssigned  = errors.New("reviewer is not assigned to this PR")
)

func isUnique(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
