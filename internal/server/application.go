package server

import (
	"database/sql"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/db/sqlite"
	"github.com/maximekuhn/warden/internal/permissions"
	"github.com/maximekuhn/warden/internal/transaction"
)

type application struct {
	authService *auth.AuthService
	permService *permissions.PermissionsService
	uowProvider transaction.UnitOfWorkProvider
}

func newApplication(db *sql.DB) application {
	authBackend := sqlite.NewSqliteAuthBackend(db)
	authService := auth.NewAuthService(authBackend)

	permBackend := sqlite.NewSqlitePermissionsBackend(db)
	permService := permissions.NewPermissionsService(permBackend)

	uowProvider := sqlite.NewSqlUnitOfWorkProvider(db)

	return application{
		authService: authService,
		permService: permService,
		uowProvider: uowProvider,
	}
}
