package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/permissions"
)

type UserRepository interface {
	GetAll(
		ctx context.Context,
		uow transaction.UnitOfWork,
		limit, offset uint,
	) ([]entities.User, error)

	Count(ctx context.Context, uow transaction.UnitOfWork) (uint, error)

	UpdatePlan(
		ctx context.Context,
		uow transaction.UnitOfWork,
		userID uuid.UUID,
		newPlan permissions.Plan,
	) error
}
