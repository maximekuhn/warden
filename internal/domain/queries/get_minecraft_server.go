package queries

import (
	"context"

	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/repositories"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type GetMinecraftServerQuery struct {
	ServerID valueobjects.MinecraftServerID
}

type GetMinecraftServerQueryHandler struct {
	uowProvider transaction.UnitOfWorkProvider
	serverRepo  repositories.MinecraftServerRepository
}

func NewGetMinecraftServerQueryHandler(
	uowProvider transaction.UnitOfWorkProvider,
	serverRepo repositories.MinecraftServerRepository,
) *GetMinecraftServerQueryHandler {
	return &GetMinecraftServerQueryHandler{
		uowProvider: uowProvider,
		serverRepo:  serverRepo,
	}
}

func (h *GetMinecraftServerQueryHandler) Handle(
	ctx context.Context,
	query GetMinecraftServerQuery,
) (*entities.MinecraftServer, bool, error) {
	return h.serverRepo.GetByID(ctx, h.uowProvider.Provide(), query.ServerID)
}
