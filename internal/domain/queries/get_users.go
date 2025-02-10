package queries

import (
	"context"

	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/repositories"
	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
)

type GetUsersQuery struct {
	Limit  uint
	Offest uint
}

type GetUsersQueryHandler struct {
	uowProvider transaction.UnitOfWorkProvider
	userService services.UserService
	userRepo    repositories.UserRepository
}

func NewGetUsersQueryHandler(
	uowProvider transaction.UnitOfWorkProvider,
	userService services.UserService,
	userRepo repositories.UserRepository,
) *GetUsersQueryHandler {
	return &GetUsersQueryHandler{
		uowProvider: uowProvider,
		userService: userService,
		userRepo:    userRepo,
	}
}

func (h *GetUsersQueryHandler) Handle(
	ctx context.Context,
	query GetUsersQuery,
) ([]UserDetails, uint, error) {
	uow := h.uowProvider.Provide()
	if err := uow.Begin(ctx); err != nil {
		return nil, 0, err
	}

	users, err := h.userService.GetAll(ctx, uow, query.Limit, query.Offest)
	if err != nil {
		return nil, 0, err
	}
	total, err := h.userRepo.Count(ctx, uow)

	if err := uow.Commit(); err != nil {
		return nil, 0, uow.Rollback()
	}

	return convertToUserDetail(users), total, err
}

func convertToUserDetail(users []entities.User) []UserDetails {
	res := make([]UserDetails, 0)
	for _, user := range users {
		res = append(res, UserDetails{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			Plan:      user.Plan,
		})
	}
	return res
}
