package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/logger"
	"github.com/maximekuhn/warden/internal/middlewares"
)

type LogoutHandler struct {
	logger      *slog.Logger
	service     *auth.AuthService
	uowProvider transaction.UnitOfWorkProvider
}

func NewLogoutHandler(l *slog.Logger, s *auth.AuthService, uowProvider transaction.UnitOfWorkProvider) *LogoutHandler {
	return &LogoutHandler{
		logger:      l,
		service:     s,
		uowProvider: uowProvider,
	}
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.post(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *LogoutHandler) post(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeLoggerWithRequestId(r.Context(), middlewares.RequestIdKey, h.logger)

	cookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		l.Error("failed to get cookie or no cookie found")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uow := h.uowProvider.Provide()

	if err := h.service.Logout(r.Context(), uow, *cookie); err != nil {
		// TODO: handle error accordingly
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("HX-Redirect", "/login")
}
