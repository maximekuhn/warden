package services

import (
	"context"
	"errors"
	"time"

	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/repositories"
	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type MinecraftServerStatusService struct {
	repo repositories.MinecraftServerRepository
}

func NewMinecraftServerStatusService(repo repositories.MinecraftServerRepository) *MinecraftServerStatusService {
	return &MinecraftServerStatusService{
		repo: repo,
	}
}

func (s *MinecraftServerStatusService) UpdateStatus(
	ctx context.Context,
	uow transaction.UnitOfWork,
	serverID valueobjects.MinecraftServerID,
	status valueobjects.MinecraftServerStatus,
) error {
	srv, found, err := s.repo.GetByID(ctx, uow, serverID)
	if err != nil {
		return err
	}
	if !found {
		return services.ErrServerNotFound
	}

	if status == srv.Status {
		return errors.New("status not changed")
	}

	now := time.Now()
	new := entities.NewMinecraftServer(
		srv.ID,
		srv.OwnerID,
		srv.Members,
		srv.Name,
		status,
		srv.CreatedAt,
		now,
	)
	return s.repo.Update(ctx, uow, *srv, *new)
}
