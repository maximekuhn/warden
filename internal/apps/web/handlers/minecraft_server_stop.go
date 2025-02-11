package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/apps/web/middlewares"
	"github.com/maximekuhn/warden/internal/apps/web/ui/components/section"
	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/commands"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/logger"
	"github.com/maximekuhn/warden/internal/permissions"
)

type MinecraftServerStopHandler struct {
	logger               *slog.Logger
	ps                   *permissions.PermissionsService
	uowProvider          transaction.UnitOfWorkProvider
	stopServerCmdHandler *commands.StopMinecraftServerCommandHandler
}

func NewMinecraftServerStopHandler(
	l *slog.Logger,
	ps *permissions.PermissionsService,
	uowProvider transaction.UnitOfWorkProvider,
	stopServerCmdHandler *commands.StopMinecraftServerCommandHandler,
) *MinecraftServerStopHandler {
	return &MinecraftServerStopHandler{
		logger:               l,
		ps:                   ps,
		uowProvider:          uowProvider,
		stopServerCmdHandler: stopServerCmdHandler,
	}
}

func (h *MinecraftServerStopHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.post(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *MinecraftServerStopHandler) post(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeLoggerWithRequestId(r.Context(), middlewares.RequestIdKey, h.logger)
	l = logger.UpgradeLoggerWithUserId(r.Context(), middlewares.LoggedUserKey, l)
	loggedUser, ok := r.Context().Value(middlewares.LoggedUserKey).(auth.User)
	if !ok {
		l.Error("logged user not found in request context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	serverIdStr := r.PathValue("serverId")
	serverId, err := valueobjects.NewMinecraftServerIDFromString(serverIdStr)
	if err != nil {
		// TODO: return proper error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check user permissions within that server
	hasPerm, err := h.ps.HasMinecraftServerPermission(
		r.Context(),
		h.uowProvider.Provide(),
		&loggedUser,
		serverId.Value(),
		permissions.ActionStopServer,
	)
	if err != nil {
		l.Error(
			"could not check user permissions",
			slog.String("errMsg", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !hasPerm {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// command
	if err := h.stopServerCmdHandler.Handle(
		r.Context(),
		commands.StopMinecraftServerCommand{
			ServerID: serverId,
		},
	); err != nil {
		l.Error(
			"failed to stop server",
			slog.String("errMsg", err.Error()),
		)
		// TODO: correct error handling
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := section.MinecraftServerStatus(
		serverId,
		valueobjects.MinecraftServerStatusStopped,
	).Render(r.Context(), w); err != nil {
		l.Error(
			"failed to render MinecraftServerStatus",
			slog.String("errMsg", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	l.Info("stop server handled successfully")
}
