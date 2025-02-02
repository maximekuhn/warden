package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/apps/web/handlers/handlerutils"
	"github.com/maximekuhn/warden/internal/apps/web/middlewares"
	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/commands"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/logger"
	"github.com/maximekuhn/warden/internal/permissions"
)

type MinecraftServerHandler struct {
	logger                          *slog.Logger
	ps                              *permissions.PermissionsService
	uowProvider                     transaction.UnitOfWorkProvider
	createMinecraftServerCmdHandler *commands.CreateMinecraftServerCommandHandler
}

func NewMinecraftServerHandler(
	l *slog.Logger,
	ps *permissions.PermissionsService,
	uowProvider transaction.UnitOfWorkProvider,
	createMinecraftServerCmdHandler *commands.CreateMinecraftServerCommandHandler,
) *MinecraftServerHandler {
	return &MinecraftServerHandler{
		logger:                          l,
		ps:                              ps,
		uowProvider:                     uowProvider,
		createMinecraftServerCmdHandler: createMinecraftServerCmdHandler,
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
	// upgrade logger and retrieve logged user from request context
	l := logger.UpgradeLoggerWithRequestId(r.Context(), middlewares.RequestIdKey, h.logger)
	l = logger.UpgradeLoggerWithUserId(r.Context(), middlewares.LoggedUserKey, l)
	loggedUser, ok := r.Context().Value(middlewares.LoggedUserKey).(auth.User)
	if !ok {
		l.Error("logged user not found in request context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// parse form
	// TODO: handle potential error
	_ = r.ParseForm()
	serverName, err := handlerutils.ToMinecraftServerNameOrReturnErrorBox(w, r.Form.Get("server-name"))
	if err != nil {
		l.Info("invalid server name", slog.String("reason", err.Error()))
		return
	}

	// check user policy
	hasPerm, err := h.ps.HasPolicy(
		r.Context(),
		h.uowProvider.Provide(),
		&loggedUser,
		permissions.PolicyCreateServer,
	)
	if err != nil {
		// TODO
		l.Error(
			"failed to retrieve logged user policy",
			slog.String("errMsg", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !hasPerm {
		// TODO
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// create minecraft server
	if err := h.createMinecraftServerCmdHandler.Handle(r.Context(), commands.CreateMinecraftServerCommand{
		Name:  serverName,
		Owner: loggedUser.ID,
	}); err != nil {
		l.Error(
			"failed to create minecraft server",
			slog.String("errMsg", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
