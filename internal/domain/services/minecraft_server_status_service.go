package services

import (
	"context"

	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

// MinecraftServerStatusService is a simple service to access the current status of a minecraft server.
type MinecraftServerStatusService interface {
	UpdateStatus(
		ctx context.Context,
		uow transaction.UnitOfWork,
		serverID valueobjects.MinecraftServerID,
		status valueobjects.MinecraftServerStatus,
	) error
}
