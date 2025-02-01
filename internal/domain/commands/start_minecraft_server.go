package commands

import (
	"context"

	"github.com/maximekuhn/warden/internal/domain/async"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type StartMinecraftServerCommand struct {
	ServerID valueobjects.MinecraftServerID
}

type StartMinecraftServerCommandHandler struct {
	eventBus    async.EventBus
	uowProvider transaction.UnitOfWorkProvider
}

func NewStartMinecraftServerCommandHandler(
	eventBus async.EventBus,
	uowProvider transaction.UnitOfWorkProvider,
) *StartMinecraftServerCommandHandler {
	return &StartMinecraftServerCommandHandler{
		eventBus:    eventBus,
		uowProvider: uowProvider,
	}
}

func (h *StartMinecraftServerCommandHandler) Handle(
	ctx context.Context,
	uow transaction.UnitOfWork,
	cmd StartMinecraftServerCommand,
) error {
	return h.eventBus.PublishStartServerEvent(async.StartServerEvent{
		ServerID: cmd.ServerID,
	})
}
