package server

import (
	"database/sql"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/commands"
	"github.com/maximekuhn/warden/internal/domain/queries"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/infra/db/sqlite"
	"github.com/maximekuhn/warden/internal/permissions"
)

type application struct {
	authService *auth.AuthService
	permService *permissions.PermissionsService
	uowProvider transaction.UnitOfWorkProvider

	createUserCmdHandler *commands.CreateUserCommandHandler

	getUserPlanQueryHandler *queries.GetUserPlanQueryHandler
}

func newApplication(db *sql.DB) application {
	authBackend := sqlite.NewSqliteAuthBackend(db)
	authService := auth.NewAuthService(authBackend)

	permBackend := sqlite.NewSqlitePermissionsBackend(db)
	permService := permissions.NewPermissionsService(permBackend)

	uowProvider := sqlite.NewSqlUnitOfWorkProvider(db)

	createUserCmdHandler := commands.NewCreateUserCommandHandler(authService, permService, uowProvider)

	getUserPlanQueryHandler := queries.NewGetUserPlanQueryHandler(permService, uowProvider)

	return application{
		authService:             authService,
		permService:             permService,
		uowProvider:             uowProvider,
		createUserCmdHandler:    createUserCmdHandler,
		getUserPlanQueryHandler: getUserPlanQueryHandler,
	}
}
