package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/valueobjects"
)

// SqliteAuthBackend implements auth.Backend
type SqliteAuthBackend struct {
	db *sql.DB
}

func NewSqliteAuthBackend(db *sql.DB) *SqliteAuthBackend {
	return &SqliteAuthBackend{db: db}
}

func (s *SqliteAuthBackend) Save(ctx context.Context, user auth.User) error {
	return errors.New("not implemented")
}

func (s *SqliteAuthBackend) GetByEmail(ctx context.Context, email valueobjects.Email) (*auth.User, bool, error) {
	return nil, false, errors.New("not implemented")
}

func (s *SqliteAuthBackend) GetBySessionId(ctx context.Context, sessionId string) (*auth.User, bool, error) {
	return nil, false, errors.New("not implemented")
}

func (s *SqliteAuthBackend) Update(ctx context.Context, old, new auth.User) error {
	return errors.New("not implemented")
}
