package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/transaction"
)

type MinecraftServerRepository interface {
	Save(
		ctx context.Context,
		uow transaction.UnitOfWork,
		ms entities.MinecraftServer,
	) error

	// GetAllForUser returns the list of all minecraft servers containing
	// the provided user as owner, admin, ...
	GetAllForUser(
		ctx context.Context,
		uow transaction.UnitOfWork,
		userID uuid.UUID,
	) ([]entities.MinecraftServer, error)

	Update(
		ctx context.Context,
		uow transaction.UnitOfWork,
		old, new entities.MinecraftServer,
	) error
}
