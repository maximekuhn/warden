package async

import (
	"context"
	"log/slog"

	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type StopServerEventListener struct {
	logger                *slog.Logger
	uowProvider           transaction.UnitOfWorkProvider
	cms                   services.ContainerManagementService
	mcServerStatusService services.MinecraftServerStatusService
}

func NewStopServerEventListener(
	logger *slog.Logger,
	uowProvider transaction.UnitOfWorkProvider,
	cms services.ContainerManagementService,
	mcServerStatusService services.MinecraftServerStatusService,

) *StopServerEventListener {
	return &StopServerEventListener{
		logger:                logger,
		uowProvider:           uowProvider,
		cms:                   cms,
		mcServerStatusService: mcServerStatusService,
	}
}

func (l *StopServerEventListener) Execute(evt StopServerEvent) {
	logger := l.upgradeLogger(evt.ServerID)
	logger.Info("received StopServerEvent")

	uow := l.uowProvider.Provide()
	ctx := context.TODO()
	if err := uow.Begin(ctx); err != nil {
		logger.Error(
			"failed to start tx",
			slog.String("errMsg", err.Error()),
		)
		return
	}

	if err := l.cms.StopMinecraftServer(ctx, uow, evt.ServerID); err != nil {
		logger.Error(
			"failed to stop minecraft server container",
			slog.String("errMsg", err.Error()),
		)
		_ = uow.Rollback()
		return
	}

	if err := l.mcServerStatusService.UpdateStatus(
		ctx,
		uow,
		evt.ServerID,
		valueobjects.MinecraftServerStatusStopped,
	); err != nil {
		logger.Error(
			"failed to update status",
			slog.String("errMsg", err.Error()),
		)
		_ = uow.Rollback()
		return

	}

	if err := uow.Commit(); err != nil {
		logger.Error(
			"failed to commit tx",
			slog.String("errMsg", err.Error()),
		)
		_ = uow.Rollback()
		return
	}

	logger.Info("successfully stopped minecraft server")
}

func (l *StopServerEventListener) upgradeLogger(serverID valueobjects.MinecraftServerID) *slog.Logger {
	return l.logger.With(slog.String("serverId", serverID.Value().String()))
}
