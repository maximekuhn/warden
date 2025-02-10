package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/apps/web/middlewares"
	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/commands"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/logger"
	"github.com/maximekuhn/warden/internal/permissions"
)

type MinecraftServerHandler struct {
	logger                  *slog.Logger
	startMcServerCmdHandler *commands.StartMinecraftServerCommandHandler
	uowProvider             transaction.UnitOfWorkProvider
	ps                      *permissions.PermissionsService
}

func NewMinecraftServerHandler(
	l *slog.Logger,
	startMcServerCmdHandler *commands.StartMinecraftServerCommandHandler,
	uowProvider transaction.UnitOfWorkProvider,
	ps *permissions.PermissionsService,
) *MinecraftServerHandler {
	return &MinecraftServerHandler{
		logger:                  l,
		startMcServerCmdHandler: startMcServerCmdHandler,
		uowProvider:             uowProvider,
		ps:                      ps,
	}
}

func (h *MinecraftServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.post(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *MinecraftServerHandler) post(w http.ResponseWriter, r *http.Request) {
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
		permissions.ActionStartServer,
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

	// start minecraft server
	err = h.startMcServerCmdHandler.Handle(r.Context(), commands.StartMinecraftServerCommand{
		ServerID: serverId,
	})
	if err != nil {
		// TODO: handle error properly
		l.Error(
			"failed to start minecraft server",
			slog.String("errMs", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	l.Info("successfully started minecraft server")
}
