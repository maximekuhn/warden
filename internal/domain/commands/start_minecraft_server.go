package commands

import (
	"context"

	"github.com/maximekuhn/warden/internal/domain/async"
	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type StartMinecraftServerCommand struct {
	ServerID valueobjects.MinecraftServerID
}

type StartMinecraftServerCommandHandler struct {
	eventBus      async.EventBus
	uowProvider   transaction.UnitOfWorkProvider
	statusService services.MinecraftServerStatusService
}

func NewStartMinecraftServerCommandHandler(
	eventBus async.EventBus,
	uowProvider transaction.UnitOfWorkProvider,
	statusService services.MinecraftServerStatusService,
) *StartMinecraftServerCommandHandler {
	return &StartMinecraftServerCommandHandler{
		eventBus:      eventBus,
		uowProvider:   uowProvider,
		statusService: statusService,
	}
}

func (h *StartMinecraftServerCommandHandler) Handle(
	ctx context.Context,
	cmd StartMinecraftServerCommand,
) error {
	// TODO: check if server is not already running or starting
	uow := h.uowProvider.Provide()
	if err := uow.Begin(ctx); err != nil {
		return err
	}

	if err := h.statusService.UpdateStatus(
		ctx,
		uow,
		cmd.ServerID,
		valueobjects.MinecraftServerStatusStarting,
	); err != nil {
		_ = uow.Rollback()
		return err
	}

	if err := uow.Commit(); err != nil {
		return uow.Rollback()
	}

	return h.eventBus.PublishStartServerEvent(async.StartServerEvent{
		ServerID: cmd.ServerID,
	})
}
