package postgres

import (
	"errors"

	apperrors "github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/apperror"
	"github.com/jackc/pgx/v5/pgconn"
)

func mapDBError(err error) error {
	if err == nil {
		return nil
	}

	// PostgreSQL specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return apperrors.Wrap(
				apperrors.Conflict,
				"resource already exists",
				err,
			)
		case "23503": // foreign_key_violation
			return apperrors.Wrap(
				apperrors.InvalidInput,
				"invalid reference",
				err,
			)
		case "23502": // not_null_violation
			return apperrors.Wrap(
				apperrors.InvalidInput,
				"missing required field",
				err,
			)
		case "23514": // check_violation
			return apperrors.Wrap(
				apperrors.InvalidInput,
				"constraint violation",
				err,
			)
		}
	}

	// Default: internal DB error
	return apperrors.Wrap(
		apperrors.Internal,
		"database operation failed",
		err,
	)
}
