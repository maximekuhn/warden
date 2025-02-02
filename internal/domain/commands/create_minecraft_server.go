package commands

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/repositories"
	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/permissions"
)

type CreateMinecraftServerCommand struct {
	Name  valueobjects.MinecraftServerName
	Owner uuid.UUID
}

type CreateMinecraftServerCommandHandler struct {
	portAllocator    services.PortAllocatorService
	userService      services.UserService
	serverRepository repositories.MinecraftServerRepository
	uowProvider      transaction.UnitOfWorkProvider
}

func NewCreateMinecraftServerCommandHandler(
	portAllocator services.PortAllocatorService,
	userService services.UserService,
	serverRepository repositories.MinecraftServerRepository,
	uowProvider transaction.UnitOfWorkProvider,
) *CreateMinecraftServerCommandHandler {
	return &CreateMinecraftServerCommandHandler{
		portAllocator:    portAllocator,
		userService:      userService,
		serverRepository: serverRepository,
		uowProvider:      uowProvider,
	}
}

func (h *CreateMinecraftServerCommandHandler) Handle(
	ctx context.Context,
	cmd CreateMinecraftServerCommand,
) error {
	uow := h.uowProvider.Provide()
	if err := uow.Begin(ctx); err != nil {
		return err
	}

	// create a new server ID and try to allocate a port
	serverID := valueobjects.GenerateMinecraftServerID()
	_, err := h.portAllocator.AllocatePort(ctx, uow, serverID)
	if err != nil {
		return err
	}

	// persist server
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

	// update user roles to include owner as admin for the newly created server
	if err := h.userService.AddRoleInServer(
		ctx,
		uow,
		cmd.Owner,
		serverID,
		permissions.RoleAdmin,
	); err != nil {
		return err
	}

	if err := uow.Commit(); err != nil {
		return uow.Rollback()
	}

	return nil
}
