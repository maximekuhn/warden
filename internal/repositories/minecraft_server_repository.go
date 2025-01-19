package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/entities"
)

type MinecraftServerRepository interface {
	Save(ctx context.Context, ms entities.MinecraftServer) error
	GetById(ctx context.Context, serverID uuid.UUID) (*entities.MinecraftServer, error)

	// GetAllForUser returns all minecraft server where user is owner or member.
	GetAllForUser(ctx context.Context, userID uuid.UUID) ([]entities.MinecraftServer, error)
}
