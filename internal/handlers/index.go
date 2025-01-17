package handlers

import (
	"log/slog"
	"net/http"

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
	if err := pages.Index().Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
