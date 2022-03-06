package postgres

import (
	"database/sql"
	"errors"
)

func EmptyOrError(err error, errorMessage string) error {
	if err == sql.ErrNoRows {
		return nil
	}
	return errors.New(errorMessage)
}
