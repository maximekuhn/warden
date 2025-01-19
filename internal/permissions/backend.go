package permissions

import (
	"context"

	"github.com/google/uuid"
)

// Backend represents whatever is used to store user permissions and related
// data. It can be a database, a file, a cache, etc...
type Backend interface {
	Save(ctx context.Context, user User) error
	Update(ctx context.Context, old, new User) error
	GetById(ctx context.Context, userID uuid.UUID) (*User, bool, error)
}
