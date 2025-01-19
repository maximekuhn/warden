package handlers

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/permissions"
)

type CreateServerHandler struct {
	logger *slog.Logger
	_      *permissions.PermissionsService
}

func NewCreateServerHandler(l *slog.Logger) *CreateServerHandler {
	return &CreateServerHandler{logger: l}
}

func (h *CreateServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *CreateServerHandler) get(w http.ResponseWriter, r *http.Request) {

}
