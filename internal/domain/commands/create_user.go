package commands

import (
	"context"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/permissions"
	"github.com/maximekuhn/warden/internal/transaction"
)

type CreateUserCommand struct {
	Email           valueobjects.Email
	Password        valueobjects.Password
	PasswordConfirm valueobjects.Password
}

type CreateUserCommandHandler struct {
	authService *auth.AuthService
	permService *permissions.PermissionsService
	uowProvider transaction.UnitOfWorkProvider
}

func NewCreateUserCommandHandler(
	as *auth.AuthService,
	ps *permissions.PermissionsService,
	uowProvider transaction.UnitOfWorkProvider,
) *CreateUserCommandHandler {
	return &CreateUserCommandHandler{
		authService: as,
		permService: ps,
		uowProvider: uowProvider,
	}
}

func (h *CreateUserCommandHandler) Handle(
	ctx context.Context,
	cmd CreateUserCommand,
) error {
	uow := h.uowProvider.Provide()
	if err := uow.Begin(ctx); err != nil {
		return err
	}

	// password and confirmation must match
	if cmd.Password != cmd.PasswordConfirm {
		return ErrPasswordAndConfirmationDontMatch
	}

	// create a auth user, which will contain session data
	userID, err := h.authService.Register(ctx, uow, cmd.Email, cmd.Password)
	if err != nil {
		return err
	}

	// create a permission user with the default subscription plan (Free)
	if err := h.permService.Create(ctx, uow, userID, permissions.PlanFree); err != nil {
		return err
	}

	if err := uow.Commit(); err != nil {
		return uow.Rollback()
	}
	return nil
}
