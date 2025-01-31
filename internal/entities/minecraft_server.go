package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/valueobjects"
)

type MinecraftServer struct {
	ID         uuid.UUID
	OwnerID    uuid.UUID   // warden account
	Members    []uuid.UUID // warden account
	Name       valueobjects.MinecraftServerName
	Status     valueobjects.MinecraftServerStatus
	HostConfig valueobjects.MinecraftServerHostConfig
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewMinecraftServer(
	id uuid.UUID,
	ownerID uuid.UUID,
	members []uuid.UUID,
	name valueobjects.MinecraftServerName,
	status valueobjects.MinecraftServerStatus,
	hostConfig valueobjects.MinecraftServerHostConfig,
	createdAt time.Time,
	updatedAt time.Time,

) *MinecraftServer {
	return &MinecraftServer{
		ID:         id,
		OwnerID:    ownerID,
		Members:    members,
		Name:       name,
		Status:     status,
		HostConfig: hostConfig,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}
