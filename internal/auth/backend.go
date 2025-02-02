package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

// Backend represents whatever is used to store users information.
// It can be a database, a file system, an in-memory cache, ...
type Backend interface {
	Save(
		ctx context.Context,
		uow transaction.UnitOfWork,
		user User,
	) error

	GetByEmail(
		ctx context.Context,
		uow transaction.UnitOfWork,
		email valueobjects.Email,
	) (*User, bool, error)

	GetByUserID(
		ctx context.Context,
		uow transaction.UnitOfWork,
		userID uuid.UUID,
	) (*User, bool, error)

	GetBySessionId(
		ctx context.Context,
		uow transaction.UnitOfWork,
		sessionId string,
	) (*User, bool, error)

	Update(
		ctx context.Context,
		uow transaction.UnitOfWork,
		old, new User,
	) error
}
