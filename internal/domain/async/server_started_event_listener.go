package async

import (
	"context"
	"log/slog"

	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type ServerStartedEventListener struct {
	logger      *slog.Logger
	msss        services.MinecraftServerStatusService
	uowProvider transaction.UnitOfWorkProvider
}

func NewServerStartedEventListener(
	l *slog.Logger,
	msss services.MinecraftServerStatusService,
	uowProvider transaction.UnitOfWorkProvider,
) *ServerStartedEventListener {
	return &ServerStartedEventListener{
		logger:      l,
		msss:        msss,
		uowProvider: uowProvider,
	}
}

func (l *ServerStartedEventListener) Execute(evt ServerStartedEvent) {
	logger := l.upgradeLogger(evt.ServerID)
	logger.Info("received server started event")

	uow := l.uowProvider.Provide()
	ctx := context.TODO()
	if err := l.msss.UpdateStatus(ctx, uow, evt.ServerID, valueobjects.MinecraftServerStatusRunning); err != nil {
		logger.Error(
			"failed to update minecraft server status",
			slog.String("errMsg", err.Error()),
		)
		return
	}
	logger.Info("minecraft server status has been updated")
}

func (l *ServerStartedEventListener) upgradeLogger(serverID valueobjects.MinecraftServerID) *slog.Logger {
	return l.logger.With(slog.String("serverId", serverID.Value().String()))
}
