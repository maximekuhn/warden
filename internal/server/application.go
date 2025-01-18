package server

import (
	"database/sql"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/db/sqlite"
)

type application struct {
	authService *auth.AuthService
}

func newApplication(db *sql.DB) application {
	authBackend := sqlite.NewSqliteAuthBackend(db)
	authService := auth.NewAuthService(authBackend)

	return application{
		authService: authService,
	}
}
