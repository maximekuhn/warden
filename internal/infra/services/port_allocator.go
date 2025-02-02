package services

import (
	"context"

	"github.com/maximekuhn/warden/internal/domain/repositories"
	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

// PortAllocator is a very dummy port allocator backed by a list of open ports, provided
// when the service is created.
//
// Once all open ports are allocated, the service is no longer able to provide
// one.
//
// TODO: once all openPorts are allocated in memory, the allocator is not working anymore
type PortAllocator struct {
	repository repositories.PortRepository
	openPorts  []int16
}

func NewPortAllocator(
	repository repositories.PortRepository,
	openPorts []int16,
) *PortAllocator {
	return &PortAllocator{
		repository: repository,
		openPorts:  openPorts,
	}
}

func (pa *PortAllocator) AllocatePort(
	ctx context.Context,
	uow transaction.UnitOfWork,
	serverID valueobjects.MinecraftServerID,
) (int16, error) {
	_, found, err := pa.repository.GetByServerID(ctx, uow, serverID)
	if err != nil {
		return 0, err
	}
	if found {
		return 0, services.ErrServerAlreadyHasAllocatedPort
	}

	port, err := pa.pickFirstAvailable()
	if err != nil {
		return 0, err
	}

	err = pa.repository.Save(ctx, uow, port, serverID)
	return port, err
}

func (pa *PortAllocator) pickFirstAvailable() (int16, error) {
	const allocatedPort = -1

	idx := -1
	for i, port := range pa.openPorts {
		if port != allocatedPort {
			idx = i
		}
	}
	if idx == -1 {
		// all ports are already allocated
		return 0, services.ErrNoPortAvailable
	}

	serverPort := pa.openPorts[idx]

	// allocate port
	pa.openPorts[idx] = allocatedPort

	return serverPort, nil
}
