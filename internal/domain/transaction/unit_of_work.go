package transaction

import "context"

// UnitOfWork represents a business transaction.
// The honest documentation would say the this is to handle sql transactions :)
type UnitOfWork interface {
	Begin(context.Context) error
	Commit() error
	Rollback() error
}

type UnitOfWorkProvider interface {
	Provide() UnitOfWork
}
