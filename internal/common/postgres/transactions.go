package postgres

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUnableGetTransactionFromContext = errors.New("unable get transaction from context")
)

type transactionCtxKey string

var transactionKey = transactionCtxKey("transaction")

// BeginTx ....
func (c *Connection) BeginTx() (*sqlx.Tx, error) {
	return c.Dbx.Beginx()
}

func (c *Connection) TryWithTransaction(ctx context.Context, fn func(stmt sqlx.Ext) error) error {
	stmt := sqlx.Ext(c.Dbx)

	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if inTransaction {
		stmt = tx
	}

	return fn(stmt)
}

func (c *Connection) MustWithTransaction(ctx context.Context, fn func(stmt *sqlx.Tx) error) error {
	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if inTransaction {
		return fn(tx)
	}

	return ErrUnableGetTransactionFromContext
}
