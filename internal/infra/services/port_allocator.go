package services

import (
	"context"
	"slices"

	"github.com/maximekuhn/warden/internal/domain/repositories"
	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type PortAllocator struct {
	repository      repositories.PortRepository
	configuredPorts []uint16
}

func NewPortAllocator(
	repository repositories.PortRepository,
	configuredPOrts []uint16,
) *PortAllocator {
	return &PortAllocator{
		repository:      repository,
		configuredPorts: configuredPOrts,
	}
}

func (pa *PortAllocator) AllocatePort(
	ctx context.Context,
	uow transaction.UnitOfWork,
	serverID valueobjects.MinecraftServerID,
) (uint16, error) {
	_, found, err := pa.repository.GetByServerID(ctx, uow, serverID)
	if err != nil {
		return 0, err
	}
	if found {
		return 0, services.ErrServerAlreadyHasAllocatedPort
	}

	port, err := pa.pickFirstAvailable(ctx, uow)
	if err != nil {
		return 0, err
	}

	err = pa.repository.Save(ctx, uow, port, serverID)
	return port, err
}

func (pa *PortAllocator) pickFirstAvailable(
	ctx context.Context,
	uow transaction.UnitOfWork,
) (uint16, error) {
	takenPorts, err := pa.repository.GetAll(ctx, uow)
	if err != nil {
		return 0, err
	}

	for _, port := range pa.configuredPorts {
		if !slices.Contains(takenPorts, port) {
			return port, nil
		}
	}

	return 0, services.ErrNoPortAvailable
}
