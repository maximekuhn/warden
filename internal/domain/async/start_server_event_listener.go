package async

import (
	"context"
	"log/slog"

	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type StartServerEventListener struct {
	logger *slog.Logger
	cms    services.ContainerManagementService
}

func NewStartServerEventListener(
	l *slog.Logger,
	cms services.ContainerManagementService,
) *StartServerEventListener {
	return &StartServerEventListener{
		logger: l,
		cms:    cms,
	}
}

func (l *StartServerEventListener) Execute(evt StartServerEvent) {
	logger := l.upgradeLogger(evt.ServerID)
	logger.Info("received event to start minecraft server")

	ctx := context.TODO()
	if err := l.cms.StartMinecraftServer(ctx, evt.ServerID); err != nil {
		logger.Error(
			"failed to start minecraft server",
			slog.String("errMsg", err.Error()),
		)
		return
	}

	logger.Info("successfully started minecraft server")
}

func (l *StartServerEventListener) upgradeLogger(serverID valueobjects.MinecraftServerID) *slog.Logger {
	return l.logger.With(slog.String("serverId", serverID.Value().String()))
}
