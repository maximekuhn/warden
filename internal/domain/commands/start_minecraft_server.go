package commands

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/async"
	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
)

type StartMinecraftServerCommand struct {
	ID uuid.UUID
}

type StartMinecraftServerCommandHandler struct {
	eventBus        *async.EventBus
	metadataService *services.MinecraftServerMetadataService
	uowProvider     transaction.UnitOfWorkProvider
}

func NewStartMinecraftServerCommandHandler(
	eventBus *async.EventBus,
	metadataService *services.MinecraftServerMetadataService,
	uowProvider transaction.UnitOfWorkProvider,
) *StartMinecraftServerCommandHandler {
	return &StartMinecraftServerCommandHandler{
		eventBus:        eventBus,
		metadataService: metadataService,
		uowProvider:     uowProvider,
	}
}

func (h *StartMinecraftServerCommandHandler) Handle(
	ctx context.Context,
	uow transaction.UnitOfWork,
	cmd StartMinecraftServerCommand,
) error {
	h.eventBus.Publish(async.EventStartServer)
	return errors.New("not yet implemented")
}
