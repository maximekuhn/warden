package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/apps/web/handlers"
	"github.com/maximekuhn/warden/internal/apps/web/middlewares"
	"github.com/maximekuhn/warden/internal/domain/async"
	"github.com/maximekuhn/warden/internal/infra/queue"
	"github.com/maximekuhn/warden/internal/logger"
)

type Server struct {
	logger *slog.Logger
	app    application
	conf   *Config
}

func NewServer(l *slog.Logger, db *sql.DB, config *Config) (*Server, error) {
	app, err := newApplication(db, config, l)
	return &Server{
		logger: l,
		app:    app,
		conf:   config,
	}, err
}

func (s *Server) Start() error {
	fs := http.FileServer(http.Dir("internal/apps/web/ui/assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	reqIdMiddleware := middlewares.NewRequestIdMiddleware()
	loggerMiddleware := middlewares.NewLoggerMiddleware(s.logger.With(slog.String(logger.LoggerNameFieldName, "LoggerMiddleware")))
	sessionMiddleware := middlewares.NewSessionMiddleware(
		s.logger.With(slog.String(logger.LoggerNameFieldName, "SessionMiddleware")),
		*s.app.authService,
		s.app.uowProvider)
	adminMiddleware := middlewares.NewAdminMiddleware(s.conf.Admin.Username, s.conf.Admin.HashedPassword)

	chain := middlewares.Chain(reqIdMiddleware, loggerMiddleware)
	chainWithSession := middlewares.Chain(chain, sessionMiddleware)
	chainAdmin := middlewares.Chain(loggerMiddleware, adminMiddleware)

	indexHandler := handlers.NewIndexHandler(
		s.logger.With(slog.String(logger.LoggerNameFieldName, "IndexHandler")),
		s.app.getUserPlanQueryHandler)
	http.Handle("/", chainWithSession.Middleware(indexHandler))

	loginHandler := handlers.NewLoginHandler(
		s.logger.With(slog.String(logger.LoggerNameFieldName, "LoginHandler")),
		s.app.authService,
		s.app.uowProvider)
	http.Handle("/login", chain.Middleware(loginHandler))

	logoutHandler := handlers.NewLogoutHandler(
		s.logger.With(slog.String(logger.LoggerNameFieldName, "LogoutHandler")),
		s.app.authService,
		s.app.uowProvider)
	http.Handle("/logout", chainWithSession.Middleware(logoutHandler))

	signupHandler := handlers.NewSignupHandler(
		s.logger.With(slog.String(logger.LoggerNameFieldName, "SignupHandler")),
		s.app.createUserCmdHandler)
	http.Handle("/signup", chain.Middleware(signupHandler))

	minecraftServersHandler := handlers.NewMinecraftServersHandler(
		s.logger.With(slog.String(logger.LoggerNameFieldName, "MinecraftServersHandler")),
		s.app.permService,
		s.app.uowProvider,
		s.app.createMinecraftServerCmdHandler,
		s.app.getMinecraftServersQueryHandler,
		s.conf.MinecraftServers.Hostname,
	)
	http.Handle("/minecraft-servers", chainWithSession.Middleware(minecraftServersHandler))

	minecraftServerHandler := handlers.NewMinecraftServerHandler(
		s.logger.With(logger.LoggerNameFieldName, "MinecraftServerHandler"),
		s.app.startMinecraftServerCmdHandler,
		s.app.uowProvider,
		s.app.permService,
		s.app.getMinecraftServerQueryHandler,
	)
	http.Handle("/minecraft-servers/{serverId}", chainWithSession.Middleware(minecraftServerHandler))

	minecraftServerStatusHandler := handlers.NewMinecraftServerStatusHandler(
		s.logger.With(slog.String(logger.LoggerNameFieldName, "MinecraftServerStatusHandler")),
		s.app.permService,
		s.app.getMinecraftServerQueryHandler,
		s.app.uowProvider,
	)
	http.Handle("/minecraft-servers/{serverId}/status", chainWithSession.Middleware(minecraftServerStatusHandler))

	minecraftServerStopHandler := handlers.NewMinecraftServerStopHandler(
		s.logger.With(logger.LoggerNameFieldName, "MinecraftServerStopHandler"),
		s.app.permService,
		s.app.uowProvider,
		s.app.stopMinecraftServerCmdHandler,
	)
	http.Handle("/minecraft-servers/{serverId}/stop", chainWithSession.Middleware(minecraftServerStopHandler))

	healthHandler := handlers.NewHealthcheckHandler(s.logger.With(slog.String(logger.LoggerNameFieldName, "HealtchCheckHandler")))
	http.Handle("/healthcheck", chain.Middleware(healthHandler))

	adminHandler := handlers.NewAdminHandler(s.logger.With(slog.String(logger.LoggerNameFieldName, "AdminHandler")))
	http.Handle("/admin", chainAdmin.Middleware(adminHandler))

	usersHandler := handlers.NewUsersHandler(
		s.logger.With(slog.String(logger.LoggerNameFieldName, "UsersHandler")),
		s.app.getUsersQueryHandler,
		s.app.updateUserPlanCmdHandler,
	)
	http.Handle("/admin/users", chainAdmin.Middleware(usersHandler))

	// start async events queue
	// this might not be the best place to do it, but it will work for now
	s.startEventsQueue()

	return http.ListenAndServe(":8787", nil)
}

func (s *Server) startEventsQueue() {
	startServerListener := async.NewStartServerEventListener(
		s.logger.With(logger.LoggerNameFieldName, "StartServerEventListener"),
		s.app.uowProvider,
		s.app.containerManagementService,
		s.app.eventBus,
		s.app.minecraftServerStatusService,
	)
	serverStartedListener := async.NewServerStartedEventListener(
		s.logger.With(logger.LoggerNameFieldName, "ServerStartedEventListener"),
		s.app.minecraftServerStatusService,
		s.app.uowProvider,
	)
	stopServerListener := async.NewStopServerEventListener(
		s.logger.With(logger.LoggerNameFieldName, "StopServerListener"),
		s.app.uowProvider,
		s.app.containerManagementService,
		s.app.minecraftServerStatusService,
	)
	q := s.app.eventBus.(*queue.EventsQueue)
	q.StartListeners(startServerListener, serverStartedListener, stopServerListener)
}
