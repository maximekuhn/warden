package permissions

import (
	"context"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

// Backend represents whatever is used to store user permissions and related
// data. It can be a database, a file, a cache, etc...
type Backend interface {
	Save(ctx context.Context, uow transaction.UnitOfWork, user User) error

	AddRole(
		ctx context.Context,
		uow transaction.UnitOfWork,
		userID uuid.UUID,
		serverID valueobjects.MinecraftServerID,
		role Role,
	) error

	GetById(
		ctx context.Context,
		uow transaction.UnitOfWork,
		userID uuid.UUID,
	) (*User, bool, error)
}
