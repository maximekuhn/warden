package queries

import (
	"context"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/permissions"
)

type GetUserPlanQuery struct {
	UserID uuid.UUID
}

type GetUserPlanQueryHandler struct {
	permService *permissions.PermissionsService
	uowProvider transaction.UnitOfWorkProvider
}

func NewGetUserPlanQueryHandler(
	permService *permissions.PermissionsService,
	uowProvider transaction.UnitOfWorkProvider,
) *GetUserPlanQueryHandler {
	return &GetUserPlanQueryHandler{
		permService: permService,
		uowProvider: uowProvider,
	}
}

func (h *GetUserPlanQueryHandler) Handle(
	ctx context.Context,
	query GetUserPlanQuery,
) (permissions.Plan, error) {
	uow := h.uowProvider.Provide()
	return h.permService.GetUserPlan(ctx, uow, query.UserID)
}
