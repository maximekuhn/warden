package handlers

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/apps/web/ui/components/lists"
	"github.com/maximekuhn/warden/internal/domain/commands"
	"github.com/maximekuhn/warden/internal/domain/queries"
	"github.com/maximekuhn/warden/internal/permissions"
)

type UsersHandler struct {
	logger                   *slog.Logger
	getUsersQueryHandler     *queries.GetUsersQueryHandler
	updateUserPlanCmdHandler *commands.UpdateUserPlanCommandHandler
}

func NewUsersHandler(
	l *slog.Logger,
	getUsersQueryHandler *queries.GetUsersQueryHandler,
	updateUserPlanCmdHandler *commands.UpdateUserPlanCommandHandler,
) *UsersHandler {
	return &UsersHandler{
		logger:                   l,
		getUsersQueryHandler:     getUsersQueryHandler,
		updateUserPlanCmdHandler: updateUserPlanCmdHandler,
	}
}

func (h *UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.getList(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.post(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *UsersHandler) getList(w http.ResponseWriter, r *http.Request) {
	query := parseGetQuery(r)
	users, _, err := h.getUsersQueryHandler.Handle(r.Context(), query)
	if err != nil {
		h.logger.Error(
			"failed to get users",
			slog.String("errMsg", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := lists.UsersList(users).Render(r.Context(), w); err != nil {
		h.logger.Error(
			"failed to render UsersList",
			slog.String("errMsg", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func parseGetQuery(_ *http.Request) queries.GetUsersQuery {
	// TODO
	return queries.GetUsersQuery{
		Limit:  100,
		Offest: 0,
	}
}

func (h *UsersHandler) post(w http.ResponseWriter, r *http.Request) {
	cmd, err := parsePostQuery(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.updateUserPlanCmdHandler.Handle(r.Context(), cmd); err != nil {
		// TODO: handle error properly
		h.logger.Error(
			"failed to update user plan",
			slog.String("errMsg", err.Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Info("successfully updated user plan")
}

func parsePostQuery(r *http.Request) (commands.UpdateUserPlanCommand, error) {
	if err := r.ParseForm(); err != nil {
		return commands.UpdateUserPlanCommand{}, err
	}

	userIdStr := r.Form.Get("userID")
	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		return commands.UpdateUserPlanCommand{}, err
	}
	planstr := r.Form.Get("newPlan")
	plan, err := permissions.PlanFromString(planstr)
	if err != nil {
		return commands.UpdateUserPlanCommand{}, err
	}

	return commands.UpdateUserPlanCommand{
		UserID:  userID,
		NewPlan: plan,
	}, nil
}
