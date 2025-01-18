package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/valueobjects"
)

// User represents data used to authenticate a user.
type User struct {
	ID uuid.UUID

	// general information
	Email         valueobjects.Email
	HashedPassord HashedPassword
	CreatedAt     time.Time

	// session
	SessionId         string    // empty if the user is logged out
	SessionExpireDate time.Time // set to UNIX epoch if the user is logged out
}

func NewUser(
	id uuid.UUID,
	email valueobjects.Email,
	hashedPassord HashedPassword,
	createdAt time.Time,
	sessionId string,
	sessionExpireDate time.Time) *User {
	return &User{
		ID:                id,
		Email:             email,
		HashedPassord:     hashedPassord,
		CreatedAt:         createdAt.UTC(),
		SessionId:         sessionId,
		SessionExpireDate: sessionExpireDate.UTC(),
	}
}

func (u User) IsLoggedOut() bool {
	return u.SessionId == ""
}
