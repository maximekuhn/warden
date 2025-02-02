package services

import (
	"context"

	"github.com/google/uuid"
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
}
