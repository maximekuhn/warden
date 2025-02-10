package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/apps/web/handlers"
	"github.com/maximekuhn/warden/internal/apps/web/middlewares"
	"github.com/maximekuhn/warden/internal/domain/async"
	"github.com/maximekuhn/warden/internal/infra/queue"
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
	loggerMiddleware := middlewares.NewLoggerMiddleware(s.logger.With(slog.Bool("LoggerMiddleware", true)))
	sessionMiddleware := middlewares.NewSessionMiddleware(
		s.logger.With(slog.Bool("SessionMiddleware", true)),
		*s.app.authService,
		s.app.uowProvider)

	chain := middlewares.Chain(reqIdMiddleware, loggerMiddleware)
	chainWithSession := middlewares.Chain(chain, sessionMiddleware)

	indexHandler := handlers.NewIndexHandler(
		s.logger.With(slog.String("handler", "IndexHandler")),
		s.app.getUserPlanQueryHandler)
	http.Handle("/", chainWithSession.Middleware(indexHandler))

	loginHandler := handlers.NewLoginHandler(
		s.logger.With(slog.String("handler", "LoginHandler")),
		s.app.authService,
		s.app.uowProvider)
	http.Handle("/login", chain.Middleware(loginHandler))

	logoutHandler := handlers.NewLogoutHandler(
		s.logger.With(slog.String("handler", "LogoutHandler")),
		s.app.authService,
		s.app.uowProvider)
	http.Handle("/logout", chainWithSession.Middleware(logoutHandler))

	signupHandler := handlers.NewSignupHandler(
		s.logger.With(slog.String("handler", "SignupHandler")),
		s.app.createUserCmdHandler)
	http.Handle("/signup", chain.Middleware(signupHandler))

	minecraftServersHandler := handlers.NewMinecraftServersHandler(
		s.logger.With(slog.String("handler", "MinecraftServersHandler")),
		s.app.permService,
		s.app.uowProvider,
		s.app.createMinecraftServerCmdHandler,
		s.app.getMinecraftServersQueryHandler,
		s.conf.MinecraftServers.Hostname,
	)
	http.Handle("/minecraft-servers", chainWithSession.Middleware(minecraftServersHandler))

	minecraftServerHandler := handlers.NewMinecraftServerHandler(
		s.logger.With("handler", "MinecraftServerHandler"),
		s.app.startMinecraftServerCmdHandler,
		s.app.uowProvider,
		s.app.permService,
	)
	http.Handle("/minecraft-servers/{serverId}", chainWithSession.Middleware(minecraftServerHandler))

	healthHandler := handlers.NewHealthcheckHandler(s.logger.With(slog.String("handler", "HealtchCheckHandler")))
	http.Handle("/healthcheck", chain.Middleware(healthHandler))

	// start async events queue
	// this might not be the best place to do it, but it will work for now
	s.startEventsQueue()

	return http.ListenAndServe(":8787", nil)
}

func (s *Server) startEventsQueue() {
	startServerListener := async.NewStartServerEventListener(
		s.logger.With("listener", "StartServerEventListener"),
		s.app.uowProvider,
		s.app.containerManagementService,
		s.app.eventBus,
	)
	serverStartedListener := async.NewServerStartedEventListener(
		s.logger.With("listener", "ServerStartedEventListener"),
		s.app.minecraftServerStatusService,
		s.app.uowProvider,
	)
	q := s.app.eventBus.(*queue.EventsQueue)
	q.StartListeners(startServerListener, serverStartedListener)
}
