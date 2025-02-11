package repositories

import (
	"context"

	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type PortRepository interface {
	Save(
		ctx context.Context,
		uow transaction.UnitOfWork,
		port uint16,
		serverID valueobjects.MinecraftServerID,
	) error

	GetByServerID(
		ctx context.Context,
		uow transaction.UnitOfWork,
		serverID valueobjects.MinecraftServerID,
	) (uint16, bool, error)

	GetAll(ctx context.Context, uow transaction.UnitOfWork) ([]uint16, error)
}
