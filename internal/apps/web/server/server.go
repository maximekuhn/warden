package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/apps/web/handlers"
	"github.com/maximekuhn/warden/internal/apps/web/middlewares"
)

type Server struct {
	logger *slog.Logger
	app    application
}

func NewServer(l *slog.Logger, db *sql.DB) *Server {
	app := newApplication(db)
	return &Server{logger: l, app: app}
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

	minecraftServerHandler := handlers.NewMinecraftServerHandler(
		s.logger.With(slog.String("handler", "MinecraftServerHandler")),
		s.app.permService,
		s.app.uowProvider,
		s.app.createMinecraftServerCmdHandler,
		s.app.getMinecraftServersQueryHandler,
	)
	http.Handle("/minecraft-servers", chainWithSession.Middleware(minecraftServerHandler))

	healthHandler := handlers.NewHealthcheckHandler(s.logger.With(slog.String("handler", "HealtchCheckHandler")))
	http.Handle("/healthcheck", chain.Middleware(healthHandler))

	return http.ListenAndServe(":8787", nil)
}
