package async

import (
	"context"
	"log/slog"

	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type StartServerEventListener struct {
	logger                *slog.Logger
	uowProvider           transaction.UnitOfWorkProvider
	cms                   services.ContainerManagementService
	eventBus              EventBus
	mcServerStatusService services.MinecraftServerStatusService
}

func NewStartServerEventListener(
	l *slog.Logger,
	uowProvider transaction.UnitOfWorkProvider,
	cms services.ContainerManagementService,
	eventBus EventBus,
	mcServerStatusService services.MinecraftServerStatusService,
) *StartServerEventListener {
	return &StartServerEventListener{
		logger:                l,
		uowProvider:           uowProvider,
		cms:                   cms,
		eventBus:              eventBus,
		mcServerStatusService: mcServerStatusService,
	}
}

func (l *StartServerEventListener) Execute(evt StartServerEvent) {
	logger := l.upgradeLogger(evt.ServerID)
	logger.Info("received StartServerEvent")

	ctx := context.TODO()
	uow := l.uowProvider.Provide()
	if err := l.cms.StartMinecraftServer(ctx, uow, evt.ServerID); err != nil {
		logger.Error(
			"failed to start minecraft server",
			slog.String("errMsg", err.Error()),
		)
		l.setServerStatusToStopped(ctx, uow, evt.ServerID, logger)
		return
	}
	logger.Info("successfully started minecraft server")

	// ignore lint for next line (false positive)
	//nolint:gosimple
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

func (l *StartServerEventListener) setServerStatusToStopped(
	ctx context.Context,
	uow transaction.UnitOfWork,
	serverID valueobjects.MinecraftServerID,
	logger *slog.Logger,
) {
	if err := l.mcServerStatusService.UpdateStatus(
		ctx,
		uow,
		serverID,
		valueobjects.MinecraftServerStatusStopped,
	); err != nil {
		logger.Error("failed to set status to stopped", slog.String("errMsg", err.Error()))
		return
	}
	logger.Info("status has been set back to stopped")
}

func (l *StartServerEventListener) upgradeLogger(serverID valueobjects.MinecraftServerID) *slog.Logger {
	return l.logger.With(slog.String("serverId", serverID.Value().String()))
}
