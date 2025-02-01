package repositories

import (
	"context"

	"github.com/maximekuhn/warden/internal/transaction"
	"github.com/maximekuhn/warden/internal/valueobjects"
)

type PortRepository interface {
	Save(
		ctx context.Context,
		uow transaction.UnitOfWork,
		port int16,
		serverID valueobjects.MinecraftServerID,
	) error

	GetByServerID(
		ctx context.Context,
		uow transaction.UnitOfWork,
		serverID valueobjects.MinecraftServerID,
	) (int16, bool, error)
}
