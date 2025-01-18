package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/ui/pages"
)

type SignupHandler struct {
	logger *slog.Logger
}

func NewSignupHandler(l *slog.Logger) *SignupHandler {
	return &SignupHandler{logger: l}
}

func (h *SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *SignupHandler) get(w http.ResponseWriter, r *http.Request) {
	if err := pages.Signup().Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
