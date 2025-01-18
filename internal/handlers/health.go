package handlers

import (
	"log/slog"
	"net/http"
)

type HealthcheckHandler struct {
	logger *slog.Logger
}

func NewHealthcheckHandler(l *slog.Logger) *HealthcheckHandler {
	return &HealthcheckHandler{logger: l}
}

func (h *HealthcheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *HealthcheckHandler) get(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte("all good, thanks!"))
}
