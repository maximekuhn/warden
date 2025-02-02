package async

import (
	"context"
	"log/slog"

	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type StartServerEventListener struct {
	logger      *slog.Logger
	uowProvider transaction.UnitOfWorkProvider
	cms         services.ContainerManagementService
	eventBus    EventBus
}

func NewStartServerEventListener(
	l *slog.Logger,
	uowProvider transaction.UnitOfWorkProvider,
	cms services.ContainerManagementService,
	eventBus EventBus,
) *StartServerEventListener {
	return &StartServerEventListener{
		logger:      l,
		uowProvider: uowProvider,
		cms:         cms,
		eventBus:    eventBus,
	}
}

func (l *StartServerEventListener) Execute(evt StartServerEvent) {
	logger := l.upgradeLogger(evt.ServerID)
	logger.Info("received event to start minecraft server")

	ctx := context.TODO()
	uow := l.uowProvider.Provide()
	if err := l.cms.StartMinecraftServer(ctx, uow, evt.ServerID); err != nil {
		logger.Error(
			"failed to start minecraft server",
			slog.String("errMsg", err.Error()),
		)
		return
	}
	logger.Info("successfully started minecraft server")

	if err := l.eventBus.PublishServerStartedEvent(ServerStartedEvent{
		ServerID: evt.ServerID,
	}); err != nil {
		logger.Error(
			"failed to publish ServerStartedEvent",
			slog.String("errMsg", err.Error()),
		)
		return
	}
	logger.Info("published ServerStartedEvent")
}

func (l *StartServerEventListener) upgradeLogger(serverID valueobjects.MinecraftServerID) *slog.Logger {
	return l.logger.With(slog.String("serverId", serverID.Value().String()))
}
