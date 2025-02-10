package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/repositories"
	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/permissions"
)

type UserService struct {
	permBackend permissions.Backend
	authBackend auth.Backend
	userRepo    repositories.UserRepository
}

func NewUserService(
	permBackend permissions.Backend,
	authBackend auth.Backend,
	userRepo repositories.UserRepository,
) *UserService {
	return &UserService{
		permBackend: permBackend,
		authBackend: authBackend,
		userRepo:    userRepo,
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

func (us *UserService) GetAll(
	ctx context.Context,
	uow transaction.UnitOfWork,
	limit, offset uint,
) ([]entities.User, error) {
	return us.userRepo.GetAll(ctx, uow, limit, offset)
}

func (us *UserService) UpdatePlan(
	ctx context.Context,
	uow transaction.UnitOfWork,
	userID uuid.UUID,
	newPlan permissions.Plan,
) error {
	u, found, err := us.permBackend.GetById(ctx, uow, userID)
	if err != nil {
		return err
	}
	if !found {
		return services.ErrUserNotFound
	}
	if u.Plan == newPlan {
		return errors.New("no need to update, plan is the same")
	}
	return us.userRepo.UpdatePlan(ctx, uow, userID, newPlan)
}
