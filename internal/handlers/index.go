package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/logger"
	"github.com/maximekuhn/warden/internal/middlewares"
	"github.com/maximekuhn/warden/internal/ui/pages"
)

type IndexHandler struct {
	logger *slog.Logger
}

func NewIndexHandler(l *slog.Logger) *IndexHandler {
	return &IndexHandler{logger: l}
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
	loggedUser, ok := r.Context().Value(middlewares.LoggedUserKey).(auth.User)
	if !ok {
		l.Error("logged user not found in request context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := pages.Index(loggedUser).Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
