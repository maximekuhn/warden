package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/ui/pages"
)

type LoginHandler struct {
	logger  *slog.Logger
	service *auth.AuthService
}

func NewLoginHandler(l *slog.Logger, s *auth.AuthService) *LoginHandler {
	return &LoginHandler{logger: l, service: s}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *LoginHandler) get(w http.ResponseWriter, r *http.Request) {
	if err := pages.Login().Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
