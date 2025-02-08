package server

import (
	"database/sql"
	"errors"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/commands"
	"github.com/maximekuhn/warden/internal/domain/queries"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/infra/db/sqlite"
	"github.com/maximekuhn/warden/internal/infra/services"
	"github.com/maximekuhn/warden/internal/permissions"
)

type application struct {
	authService *auth.AuthService
	permService *permissions.PermissionsService
	uowProvider transaction.UnitOfWorkProvider

	createUserCmdHandler            *commands.CreateUserCommandHandler
	createMinecraftServerCmdHandler *commands.CreateMinecraftServerCommandHandler

	getUserPlanQueryHandler         *queries.GetUserPlanQueryHandler
	getMinecraftServersQueryHandler *queries.GetMinecraftServersQueryHandler
}

func newApplication(db *sql.DB, conf *Config) (application, error) {
	authBackend := sqlite.NewSqliteAuthBackend(db)
	authService := auth.NewAuthService(authBackend)

	permBackend := sqlite.NewSqlitePermissionsBackend(db)
	permService := permissions.NewPermissionsService(permBackend)

	userService := services.NewUserService(permBackend, authBackend)

	portRepository := sqlite.NewSqlitePortRepository(db)
	minecraftServerRepository := sqlite.NewSqliteMinecraftServerRepository(db)
	portAllocatorService := services.NewPortAllocator(portRepository, conf.MinecraftServers.PortAllocation.Ports)

	// create docker container management service and checks if mc server image is
	// already built. If not, return an error.
	dockerContainerMngmtService, err := services.NewDockerContainerManagementService()
	if err != nil {
		return application{}, err
	}
	imageExists, err := dockerContainerMngmtService.EnsureMinecraftServerImageExists()
	if err != nil {
		return application{}, err
	}
	if !imageExists {
		return application{}, errors.New("minecraft server image not found")
	}

	uowProvider := sqlite.NewSqlUnitOfWorkProvider(db)

	// commands
	createUserCmdHandler := commands.NewCreateUserCommandHandler(authService, permService, uowProvider)
	createMinecraftServerCmdHandler := commands.NewCreateMinecraftServerCommandHandler(
		portAllocatorService,
		userService,
		minecraftServerRepository,
		uowProvider,
	)

	// queries
	getUserPlanQueryHandler := queries.NewGetUserPlanQueryHandler(permService, uowProvider)
	getMinecraftServersQueryHandler := queries.NewGetMinecraftServersQueryHandler(
		userService,
		portRepository,
		minecraftServerRepository,
		uowProvider,
	)

	return application{
		authService:                     authService,
		permService:                     permService,
		uowProvider:                     uowProvider,
		createUserCmdHandler:            createUserCmdHandler,
		createMinecraftServerCmdHandler: createMinecraftServerCmdHandler,
		getUserPlanQueryHandler:         getUserPlanQueryHandler,
		getMinecraftServersQueryHandler: getMinecraftServersQueryHandler,
	}, nil
}
