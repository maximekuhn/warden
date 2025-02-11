package commands

import (
	"context"

	"github.com/maximekuhn/warden/internal/domain/async"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type StopMinecraftServerCommand struct {
	ServerID valueobjects.MinecraftServerID
}

type StopMinecraftServerCommandHandler struct {
	uowProvider transaction.UnitOfWorkProvider
	eventBus    async.EventBus
}

func NewStopMinecraftServerCommandHandler(
	uowProvider transaction.UnitOfWorkProvider,
	eventBus async.EventBus,
) *StopMinecraftServerCommandHandler {
	return &StopMinecraftServerCommandHandler{
		uowProvider: uowProvider,
		eventBus:    eventBus,
	}
}

func (h *StopMinecraftServerCommandHandler) Handle(
	ctx context.Context,
	cmd StopMinecraftServerCommand,
) error {
	// TODO: status "Stopping" ?
	return h.eventBus.PublishStopServerEvent(async.StopServerEvent{
		ServerID: cmd.ServerID,
	})
}
