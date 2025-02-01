package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/queries"
	"github.com/maximekuhn/warden/internal/logger"
	"github.com/maximekuhn/warden/internal/middlewares"
	"github.com/maximekuhn/warden/internal/ui/pages"
)

type IndexHandler struct {
	logger                  *slog.Logger
	getUserPlanQueryHandler *queries.GetUserPlanQueryHandler
}

func NewIndexHandler(l *slog.Logger, queryHandler *queries.GetUserPlanQueryHandler) *IndexHandler {
	return &IndexHandler{logger: l, getUserPlanQueryHandler: queryHandler}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *IndexHandler) get(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeLoggerWithRequestId(r.Context(), middlewares.RequestIdKey, h.logger)
	l = logger.UpgradeLoggerWithUserId(r.Context(), middlewares.LoggedUserKey, l)

	loggedUser, ok := r.Context().Value(middlewares.LoggedUserKey).(auth.User)
	if !ok {
		l.Error("logged user not found in request context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get user plan, it is needed to correctly display the index page.
	// For instance, if the user has the plan free, we should not display a button
	// to create a new minecraft, as this plan doesn't allow him to do so.
	plan, err := h.getUserPlanQueryHandler.Handle(r.Context(), queries.GetUserPlanQuery{
		UserID: loggedUser.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		l.Error("could not retrieve user plan", slog.String("errMsg", err.Error()))
		return
	}

	if err := pages.Index(loggedUser, plan).Render(r.Context(), w); err != nil {
		l.Error("failed to render pages.Index", slog.String("errMsg", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
