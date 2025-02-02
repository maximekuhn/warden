package queries

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/repositories"
	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
)

type GetMinecraftServersQuery struct {
	UserID uuid.UUID
}

type GetMinecraftServersQueryHandler struct {
	userService services.UserService
	portRepo    repositories.PortRepository
	msRepo      repositories.MinecraftServerRepository
	uowProvider transaction.UnitOfWorkProvider
}

func NewGetMinecraftServersQueryHandler(
	userService services.UserService,
	portRepo repositories.PortRepository,
	msRepo repositories.MinecraftServerRepository,
	uowProvider transaction.UnitOfWorkProvider,
) *GetMinecraftServersQueryHandler {
	return &GetMinecraftServersQueryHandler{
		userService: userService,
		portRepo:    portRepo,
		msRepo:      msRepo,
		uowProvider: uowProvider,
	}
}

func (h *GetMinecraftServersQueryHandler) Handle(
	ctx context.Context,
	query GetMinecraftServersQuery,
) ([]MinecraftServerOverview, error) {
	uow := h.uowProvider.Provide()
	if err := uow.Begin(ctx); err != nil {
		return nil, err
	}

	servers, err := h.msRepo.GetAllForUser(ctx, uow, query.UserID)
	if err != nil {
		return nil, err
	}

	overviews := make([]MinecraftServerOverview, 0)

	// note: could be optimised in terms of query, fine for now
	for _, server := range servers {
		port, found, err := h.portRepo.GetByServerID(ctx, uow, server.ID)
		if err != nil {
			return overviews, err
		}
		if !found {
			return overviews, fmt.Errorf(
				"port not found for server %s",
				server.ID.Value().String(),
			)
		}

		role, found, err := h.userService.GetUserRoleInServer(ctx, uow, query.UserID, server.ID)
		if err != nil {
			return overviews, err
		}
		if !found {
			return overviews, fmt.Errorf(
				"role not found for user %s in server %s",
				query.UserID,
				server.ID.Value().String(),
			)
		}

		email, found, err := h.userService.GetUserEmail(ctx, uow, server.OwnerID)
		if err != nil {
			return overviews, err
		}
		if !found {
			return overviews, fmt.Errorf("email not found for user %s", query.UserID)
		}

		overview := MinecraftServerOverview{
			ID:             server.ID,
			Name:           server.Name,
			Status:         server.Status,
			Owner:          email,
			LoggedUserRole: role,
			Port:           port,
		}
		overviews = append(overviews, overview)
	}

	if err := uow.Commit(); err != nil {
		return nil, uow.Rollback()
	}
	return overviews, nil
}
