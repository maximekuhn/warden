package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/permissions"
)

type UserService interface {
	GetUserRoleInServer(
		ctx context.Context,
		uow transaction.UnitOfWork,
		userID uuid.UUID,
		serverID valueobjects.MinecraftServerID,
	) (permissions.Role, bool, error)

	GetUserEmail(
		ctx context.Context,
		uow transaction.UnitOfWork,
		userID uuid.UUID,
	) (valueobjects.Email, bool, error)

	AddRoleInServer(
		ctx context.Context,
		uow transaction.UnitOfWork,
		userID uuid.UUID,
		serverID valueobjects.MinecraftServerID,
		role permissions.Role,
	) error

	GetAll(
		ctx context.Context,
		uow transaction.UnitOfWork,
		limit, offset uint,
	) ([]entities.User, error)

	UpdatePlan(
		ctx context.Context,
		uow transaction.UnitOfWork,
		userID uuid.UUID,
		newPlan permissions.Plan,
	) error
}
