package sqlite

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/transaction"
	"github.com/maximekuhn/warden/internal/valueobjects"
)

func createTmpDb() *sql.DB {
	f, err := os.CreateTemp("", "test-db-*.sqlite3")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		panic(err)
	}
	return db
}

func createTmpDbWithAllMigrationsApplied() *sql.DB {
	db := createTmpDb()
	if err := Migrate(db); err != nil {
		db.Close()
		panic(err)
	}
	return db
}

func createUser(
	id uuid.UUID,
	emailStr string,
	createdAt time.Time,
	sessionId string,
	sessionExpiryDate time.Time,
) auth.User {
	email, err := valueobjects.NewEmail(emailStr)
	if err != nil {
		panic(err)
	}
	return *auth.NewUser(
		id,
		email,
		[]byte("hashedpasswordverysecure"),
		createdAt,
		sessionId,
		sessionExpiryDate,
	)
}

func createContextWith5MinutesTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Minute)
}

func newStringPointer(val string) *string {
	return &val
}

func newTimePointer(time time.Time) *time.Time {
	return &time
}

func createUnitOfWork(db *sql.DB) transaction.UnitOfWork {
	return NewSqlUnitOfWork(db)
}
