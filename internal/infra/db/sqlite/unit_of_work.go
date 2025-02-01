package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/maximekuhn/warden/internal/domain/transaction"
)

type SqlUnitOfWork struct {
	db *sql.DB
	tx *sql.Tx
}

func NewSqlUnitOfWork(db *sql.DB) *SqlUnitOfWork {
	return &SqlUnitOfWork{
		db: db,
		tx: nil,
	}
}

func (uow *SqlUnitOfWork) Begin(ctx context.Context) error {
	if uow.tx != nil {
		return errors.New("transaction already started")
	}
	tx, err := uow.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	uow.tx = tx
	return nil
}

func (uow *SqlUnitOfWork) Commit() error {
	if uow.tx == nil {
		return errors.New("commit: no transaction has been started")
	}
	return uow.tx.Commit()
}
func (uow *SqlUnitOfWork) Rollback() error {
	if uow.tx == nil {
		return errors.New("rollback: no transaction has been started")
	}
	return uow.tx.Rollback()
}

func (uow *SqlUnitOfWork) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if uow.tx != nil {
		return uow.tx.ExecContext(ctx, query, args...)
	}
	return uow.db.ExecContext(ctx, query, args...)
}
func (uow *SqlUnitOfWork) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if uow.tx != nil {
		return uow.tx.QueryRowContext(ctx, query, args...)
	}
	return uow.db.QueryRowContext(ctx, query, args...)
}

type SqlUnitOfWorkProvider struct {
	db *sql.DB
}

func NewSqlUnitOfWorkProvider(db *sql.DB) *SqlUnitOfWorkProvider {
	return &SqlUnitOfWorkProvider{db: db}
}

func (p *SqlUnitOfWorkProvider) Provide() transaction.UnitOfWork {
	return NewSqlUnitOfWork(p.db)
}

// castUnitOfWorkOrPanic transforms the provided transaction.UnitOfWork into a
// *SqlUnitOfWork or panics.
func castUnitOfWorkOrPanic(uow transaction.UnitOfWork) *SqlUnitOfWork {
	suow, ok := uow.(*SqlUnitOfWork)
	if !ok {
		panic("provided unit of work must be of type *SqlUnitOfWork")
	}
	return suow
}
