package auth

import (
	"context"

	"github.com/maximekuhn/warden/internal/valueobjects"
)

// Backend represents whatever is used to store users information.
// It can be a database, a file system, an in-memory cache, ...
type Backend interface {
	Save(ctx context.Context, user User) error
	GetByEmail(ctx context.Context, email valueobjects.Email) (*User, bool, error)
	GetBySessionId(ctx context.Context, sessionId string) (*User, bool, error)
	Update(ctx context.Context, old, new User) error
}
