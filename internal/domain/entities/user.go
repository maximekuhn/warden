package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/permissions"
)

type User struct {
	ID        uuid.UUID
	Email     valueobjects.Email
	Plan      permissions.Plan
	CreatedAt time.Time
}

func NewUser(
	id uuid.UUID,
	email valueobjects.Email,
	plan permissions.Plan,
	createdAt time.Time,
) *User {
	return &User{
		ID:        id,
		Email:     email,
		Plan:      plan,
		CreatedAt: createdAt,
	}
}
