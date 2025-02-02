package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/permissions"
)

type UserService struct {
	permBackend permissions.Backend
	authBackend auth.Backend
}

func NewUserService(
	permBackend permissions.Backend,
	authBackend auth.Backend,
) *UserService {
	return &UserService{
		permBackend: permBackend,
		authBackend: authBackend,
	}
}
func (us *UserService) GetUserRoleInServer(
	ctx context.Context,
	uow transaction.UnitOfWork,
	userID uuid.UUID,
	serverID valueobjects.MinecraftServerID,
) (permissions.Role, bool, error) {
	user, found, err := us.permBackend.GetById(ctx, uow, userID)
	if err != nil {
		return permissions.RoleViewer, false, err
	}
	if !found {
		return permissions.RoleViewer, false, nil
	}
	role, found := user.Roles[serverID.Value()]
	if !found {
		return permissions.RoleViewer, false, nil
	}
	return role, true, nil
}

func (us *UserService) GetUserEmail(
	ctx context.Context,
	uow transaction.UnitOfWork,
	userID uuid.UUID,
) (valueobjects.Email, bool, error) {
	user, found, err := us.authBackend.GetByUserID(ctx, uow, userID)
	if err != nil {
		return valueobjects.Email{}, false, err
	}
	if !found {
		return valueobjects.Email{}, false, nil
	}
	return user.Email, true, nil
}

func (us *UserService) AddRoleInServer(
	ctx context.Context,
	uow transaction.UnitOfWork,
	userID uuid.UUID,
	serverID valueobjects.MinecraftServerID,
	role permissions.Role,
) error {
	return us.permBackend.AddRole(ctx, uow, userID, serverID, role)
}
