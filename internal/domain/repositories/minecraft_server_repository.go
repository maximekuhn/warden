package repositories

import (
	"context"

	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/transaction"
)

type MinecraftServerRepository interface {
	Save(
		ctx context.Context,
		uow transaction.UnitOfWork,
		ms entities.MinecraftServer,
	) error
}
