package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/apps/web/ui/pages"
)

type AdminHandler struct {
	logger *slog.Logger
}

func NewAdminHandler(l *slog.Logger) *AdminHandler {
	return &AdminHandler{
		logger: l,
	}
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *AdminHandler) get(w http.ResponseWriter, r *http.Request) {
	if err := pages.Admin().Render(r.Context(), w); err != nil {
		h.logger.Error("failed to render Admin page", slog.String("errMsg", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
