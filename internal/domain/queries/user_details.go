package queries

import (
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/permissions"
)

type UserDetails struct {
	ID        uuid.UUID
	Email     valueobjects.Email
	CreatedAt time.Time
	Plan      permissions.Plan
}
