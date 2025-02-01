package services

import (
	"context"

	"github.com/maximekuhn/warden/internal/transaction"
	"github.com/maximekuhn/warden/internal/valueobjects"
)

// PortAllocatorService is a service to assign and persist ports for
// minecraft servers.
type PortAllocatorService interface {

	// AllocatePort picks an available port and assign it to the given minecraft
	// server and returns it. The port will persist across reboots.
	//
	// If no port is available, an error of type ErrNoPortAvailable is returned.
	//
	// It is not possible for the same server to have 2 allocated ports.
	AllocatePort(
		ctx context.Context,
		uow transaction.UnitOfWork,
		serverID valueobjects.MinecraftServerID,
	) (int16, error)
}
