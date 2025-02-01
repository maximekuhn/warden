package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type MinecraftServer struct {
	ID        valueobjects.MinecraftServerID
	OwnerID   uuid.UUID   // warden account
	Members   []uuid.UUID // warden account
	Name      valueobjects.MinecraftServerName
	Status    valueobjects.MinecraftServerStatus
	Port      int16
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewMinecraftServer(
	id valueobjects.MinecraftServerID,
	ownerID uuid.UUID,
	members []uuid.UUID,
	name valueobjects.MinecraftServerName,
	status valueobjects.MinecraftServerStatus,
	port int16,
	createdAt time.Time,
	updatedAt time.Time,
) *MinecraftServer {
	return &MinecraftServer{
		ID:        id,
		OwnerID:   ownerID,
		Members:   members,
		Name:      name,
		Status:    status,
		Port:      port,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
