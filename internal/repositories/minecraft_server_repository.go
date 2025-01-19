package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/entities"
)

type MinecraftServerRepository interface {
	Save(
		ctx context.Context,
		// uow transaction.UnitOfWork,
		ms entities.MinecraftServer,
	) error

	GetById(
		ctx context.Context,
		// uow transaction.UnitOfWork,
		serverID uuid.UUID,
	) (*entities.MinecraftServer, error)

	// GetAllForUser returns all minecraft server where user is owner or member.
	GetAllForUser(
		ctx context.Context,
		// uow transaction.UnitOfWork,
		userID uuid.UUID,
	) ([]entities.MinecraftServer, error)
}
