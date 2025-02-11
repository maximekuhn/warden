package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/apps/web/middlewares"
	"github.com/maximekuhn/warden/internal/apps/web/ui/components/section"
	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/queries"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/logger"
	"github.com/maximekuhn/warden/internal/permissions"
)

type MinecraftServerStatusHandler struct {
	logger                  *slog.Logger
	ps                      *permissions.PermissionsService
	getMcServerQueryHandler *queries.GetMinecraftServerQueryHandler
	uowProvider             transaction.UnitOfWorkProvider
}

func NewMinecraftServerStatusHandler(
	l *slog.Logger,
	ps *permissions.PermissionsService,
	getMcServerQueryHandler *queries.GetMinecraftServerQueryHandler,
	uowProvider transaction.UnitOfWorkProvider,
) *MinecraftServerStatusHandler {
	return &MinecraftServerStatusHandler{
		logger:                  l,
		ps:                      ps,
		getMcServerQueryHandler: getMcServerQueryHandler,
		uowProvider:             uowProvider,
	}
}

func (h *MinecraftServerStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.getStatus(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *MinecraftServerStatusHandler) getStatus(w http.ResponseWriter, r *http.Request) {
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
		permissions.ActionViewServer,
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

	// query server
	mcServer, found, err := h.getMcServerQueryHandler.Handle(
		r.Context(),
		queries.GetMinecraftServerQuery{
			ServerID: serverId,
		},
	)
	if err != nil {
		// TODO: proper error handling
		l.Error(
			"failed to retrieve minecraft server",
			slog.String("errMsg", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := section.MinecraftServerStatus(
		mcServer.ID,
		mcServer.Status,
	).Render(r.Context(), w); err != nil {
		l.Error(
			"failed to render MinecraftServerStatus",
			slog.String("errMsg", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l.Info("successfully returned status")
}
