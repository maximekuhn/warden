package services

import (
	"context"

	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

// ContainerManagementService is a service dedicated to starting and stopping
// containers. For now, containers are restricted  to minecraft server containers.
type ContainerManagementService interface {
	StartMinecraftServer(
		ctx context.Context,
		uow transaction.UnitOfWork,
		serverID valueobjects.MinecraftServerID,
	) error

	StopMinecraftServer(
		ctx context.Context,
		uow transaction.UnitOfWork,
		serverID valueobjects.MinecraftServerID,
	) error
}
