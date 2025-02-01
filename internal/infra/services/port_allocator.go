package services

import (
	"context"

	"github.com/maximekuhn/warden/internal/repositories"
	"github.com/maximekuhn/warden/internal/services"
	"github.com/maximekuhn/warden/internal/transaction"
	"github.com/maximekuhn/warden/internal/valueobjects"
)

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
		openPorts:  []int16{},
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
