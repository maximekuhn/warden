package server

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/async"
	"github.com/maximekuhn/warden/internal/domain/commands"
	"github.com/maximekuhn/warden/internal/domain/queries"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/infra/db/sqlite"
	"github.com/maximekuhn/warden/internal/infra/queue"
	"github.com/maximekuhn/warden/internal/infra/services"
	"github.com/maximekuhn/warden/internal/permissions"
)

type application struct {
	authService *auth.AuthService
	permService *permissions.PermissionsService
	uowProvider transaction.UnitOfWorkProvider

	containerManagementService   *services.DockerContainerManagementService
	minecraftServerStatusService *services.MinecraftServerStatusService

	eventBus async.EventBus

	createUserCmdHandler            *commands.CreateUserCommandHandler
	createMinecraftServerCmdHandler *commands.CreateMinecraftServerCommandHandler
	startMinecraftServerCmdHandler  *commands.StartMinecraftServerCommandHandler
	updateUserPlanCmdHandler        *commands.UpdateUserPlanCommandHandler

	getUserPlanQueryHandler         *queries.GetUserPlanQueryHandler
	getMinecraftServersQueryHandler *queries.GetMinecraftServersQueryHandler
	getUsersQueryHandler            *queries.GetUsersQueryHandler
	getMinecraftServerQueryHandler  *queries.GetMinecraftServerQueryHandler
}

func newApplication(
	db *sql.DB,
	conf *Config,
	l *slog.Logger,
) (application, error) {
	authBackend := sqlite.NewSqliteAuthBackend(db)
	authService := auth.NewAuthService(authBackend)

	permBackend := sqlite.NewSqlitePermissionsBackend(db)
	permService := permissions.NewPermissionsService(permBackend)

	userRepository := sqlite.NewSqliteUserRepository(db)
	userService := services.NewUserService(permBackend, authBackend, userRepository)

	portRepository := sqlite.NewSqlitePortRepository(db)
	minecraftServerRepository := sqlite.NewSqliteMinecraftServerRepository(db)
	portAllocatorService := services.NewPortAllocator(portRepository, conf.MinecraftServers.PortAllocation.Ports)
	minecraftServerStatusService := services.NewMinecraftServerStatusService(minecraftServerRepository)

	// create docker container management service and checks if mc server image is
	// already built. If not, return an error.
	dockerContainerMngmtService, err := services.NewDockerContainerManagementService(portRepository)
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

	eventsQueue := queue.NewEventsQeue(5, l.With(slog.Bool("EventBus", true)))

	// commands
	createUserCmdHandler := commands.NewCreateUserCommandHandler(authService, permService, uowProvider)
	createMinecraftServerCmdHandler := commands.NewCreateMinecraftServerCommandHandler(
		portAllocatorService,
		userService,
		minecraftServerRepository,
		uowProvider,
	)
	startMinecraftServerCmdHandler := commands.NewStartMinecraftServerCommandHandler(
		eventsQueue,
		uowProvider,
		minecraftServerStatusService,
	)
	updateUserPlanCmdHandler := commands.NewUpdateUserPlanCommandHandler(uowProvider, userService)

	// queries
	getUserPlanQueryHandler := queries.NewGetUserPlanQueryHandler(permService, uowProvider)
	getMinecraftServersQueryHandler := queries.NewGetMinecraftServersQueryHandler(
		userService,
		portRepository,
		minecraftServerRepository,
		uowProvider,
	)
	getUsersQueryHandler := queries.NewGetUsersQueryHandler(
		uowProvider,
		userService,
		userRepository,
	)
	getMcServerQueryHandler := queries.NewGetMinecraftServerQueryHandler(uowProvider, minecraftServerRepository)

	return application{
		authService:                     authService,
		permService:                     permService,
		uowProvider:                     uowProvider,
		containerManagementService:      dockerContainerMngmtService,
		minecraftServerStatusService:    minecraftServerStatusService,
		eventBus:                        eventsQueue,
		createUserCmdHandler:            createUserCmdHandler,
		createMinecraftServerCmdHandler: createMinecraftServerCmdHandler,
		startMinecraftServerCmdHandler:  startMinecraftServerCmdHandler,
		updateUserPlanCmdHandler:        updateUserPlanCmdHandler,
		getUserPlanQueryHandler:         getUserPlanQueryHandler,
		getMinecraftServersQueryHandler: getMinecraftServersQueryHandler,
		getUsersQueryHandler:            getUsersQueryHandler,
		getMinecraftServerQueryHandler:  getMcServerQueryHandler,
	}, nil
}
