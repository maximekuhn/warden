package services

import (
	"context"
	"errors"

	"github.com/maximekuhn/warden/internal/domain/repositories"
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
	return errors.New("not yet implemented")
}
