package commands

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/repositories"
	"github.com/maximekuhn/warden/internal/services"
	"github.com/maximekuhn/warden/internal/transaction"
)

type CreateMinecraftServerCommand struct {
	Name  valueobjects.MinecraftServerName
	Owner uuid.UUID
}

type CreateMinecraftServerCommandHandler struct {
	portAllocator    services.PortAllocatorService
	serverRepository repositories.MinecraftServerRepository
	uowProvider      transaction.UnitOfWorkProvider
}

func NewCreateMinecraftServerCommandHandler(
	portAllocator services.PortAllocatorService,
	serverRepository repositories.MinecraftServerRepository,
	uowProvider transaction.UnitOfWorkProvider,
) *CreateMinecraftServerCommandHandler {
	return &CreateMinecraftServerCommandHandler{
		portAllocator:    portAllocator,
		serverRepository: serverRepository,
		uowProvider:      uowProvider,
	}
}

func (h *CreateMinecraftServerCommandHandler) Handle(
	ctx context.Context,
	cmd CreateMinecraftServerCommand,
) error {
	uow := h.uowProvider.Provide()
	uow.Begin(ctx)

	serverID := valueobjects.NewMinecraftServerID()
	_, err := h.portAllocator.AllocatePort(ctx, uow, serverID)
	if err != nil {
		return err
	}

	now := time.Now()
	server := entities.NewMinecraftServer(
		serverID,
		cmd.Owner,
		make([]uuid.UUID, 0),
		cmd.Name,
		valueobjects.MinecraftServerStatusStopped,
		now,
		now,
	)
	if err := h.serverRepository.Save(ctx, uow, *server); err != nil {
		return err
	}

	if err := uow.Commit(); err != nil {
		return uow.Rollback()
	}

	return nil

}
