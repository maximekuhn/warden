package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
	query := `
    INSERT INTO auth (
        user_id, email, hashed_password,
        created_at, session_id, session_expire_date
    )
    VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err := s.db.ExecContext(
		ctx,
		query,
		user.ID.String(),
		user.Email.Value(),
		user.HashedPassord,
		user.CreatedAt,
		user.SessionId,
		user.SessionExpireDate)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return auth.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (s *SqliteAuthBackend) GetByEmail(ctx context.Context, email valueobjects.Email) (*auth.User, bool, error) {
	query := `
    SELECT user_id, email, hashed_password, created_at, session_id, session_expire_date
    FROM auth
    WHERE email = ?
    `
	return s.getByStringAttribute(ctx, query, email.Value())
}

func (s *SqliteAuthBackend) GetBySessionId(ctx context.Context, sessionId string) (*auth.User, bool, error) {
	query := `
    SELECT user_id, email, hashed_password, created_at, session_id, session_expire_date
    FROM auth
    WHERE session_id = ?
    `
	return s.getByStringAttribute(ctx, query, sessionId)
}

func (s *SqliteAuthBackend) Update(ctx context.Context, old, new auth.User) error {
	if old.ID != new.ID {
		return errors.New("user ID can't be updated")
	}

	updates := make([]string, 0)
	args := make([]interface{}, 0)

	if old.Email != new.Email {
		updates = append(updates, "email = ? ")
		args = append(args, new.Email.Value())
	}

	if string(old.HashedPassord) != string(new.HashedPassord) {
		updates = append(updates, "hashed_password = ? ")
		args = append(args, new.HashedPassord)
	}

	if old.SessionId != new.SessionId {
		updates = append(updates, "session_id = ? ")
		args = append(args, new.SessionId)
	}

	if old.SessionExpireDate != new.SessionExpireDate {
		updates = append(updates, "session_expire_date = ? ")
		args = append(args, new.SessionExpireDate)
	}

	if len(updates) == 0 {
		// no updates to perform
		return nil
	}

	query := fmt.Sprintf("UPDATE auth SET %s WHERE user_id = ?", strings.Join(updates, ", "))
	args = append(args, new.ID)

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("expected to affter 1 row but affected %d row(s)", affected)
	}
	return nil
}

func (s *SqliteAuthBackend) getByStringAttribute(
	ctx context.Context,
	query string,
	attributeValue string,
) (*auth.User, bool, error) {
	row := s.db.QueryRowContext(ctx, query, attributeValue)
	if err := row.Err(); err != nil {
		return nil, false, err
	}

	var id uuid.UUID
	var emailStr string
	var hashedPassword []byte
	var createdAt time.Time
	var sessionId string
	var sesionExpireDate time.Time
	if err := row.Scan(&id, &emailStr, &hashedPassword, &createdAt, &sessionId, &sesionExpireDate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, auth.ErrUserNotFound
		}
		return nil, false, err
	}

	email, err := valueobjects.NewEmail(emailStr)
	if err != nil {
		return nil, false, err
	}

	user := auth.NewUser(
		id,
		email,
		hashedPassword,
		createdAt,
		sessionId,
		sesionExpireDate,
	)
	return user, true, nil
}
