package commands

import (
	"context"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/permissions"
)

type UpdateUserPlanCommand struct {
	UserID  uuid.UUID
	NewPlan permissions.Plan
}

type UpdateUserPlanCommandHandler struct {
	uowProvider transaction.UnitOfWorkProvider
	userService services.UserService
}

func NewUpdateUserPlanCommandHandler(
	uowProvider transaction.UnitOfWorkProvider,
	userService services.UserService,
) *UpdateUserPlanCommandHandler {
	return &UpdateUserPlanCommandHandler{
		uowProvider: uowProvider,
		userService: userService,
	}
}

func (h *UpdateUserPlanCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateUserPlanCommand,
) error {
	uow := h.uowProvider.Provide()
	if err := uow.Begin(ctx); err != nil {
		return err
	}

	err := h.userService.UpdatePlan(ctx, uow, cmd.UserID, cmd.NewPlan)

	if err := uow.Commit(); err != nil {
		return uow.Rollback()
	}
	return err
}
