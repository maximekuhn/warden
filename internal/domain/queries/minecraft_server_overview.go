package queries

import (
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/permissions"
)

type MinecraftServerOverview struct {
	ID             valueobjects.MinecraftServerID
	Name           valueobjects.MinecraftServerName
	Status         valueobjects.MinecraftServerStatus
	Owner          valueobjects.Email
	LoggedUserRole permissions.Role
	Port           uint16
}
