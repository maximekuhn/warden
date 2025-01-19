package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/middlewares"
)

type LogoutHandler struct {
	logger  *slog.Logger
	service *auth.AuthService
}

func NewLogoutHandler(l *slog.Logger, s *auth.AuthService) *LogoutHandler {
	return &LogoutHandler{logger: l, service: s}
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.post(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *LogoutHandler) post(w http.ResponseWriter, r *http.Request) {
	reqId, ok := r.Context().Value(middlewares.RequestIdKey).(string)
	if !ok {
		reqId = "unknown"
	}
	l := h.logger.With(slog.String("requestId", reqId))

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		l.Error("failed to get cookie or no cookie found")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := h.service.Logout(r.Context(), *cookie); err != nil {
		// TODO: handle error accordingly
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("HX-Redirect", "/login")
}
