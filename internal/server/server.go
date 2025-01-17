package server

import (
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/handlers"
	"github.com/maximekuhn/warden/internal/middlewares"
)

type Server struct {
	logger *slog.Logger
}

func NewServer(l *slog.Logger) *Server {
	return &Server{logger: l}
}

func (s *Server) Start() error {
	fs := http.FileServer(http.Dir("internal/ui/assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	reqIdMiddleware := middlewares.NewRequestIdMiddleware()
	loggerMiddleware := middlewares.NewLoggerMiddleware(s.logger.With(slog.String("middleware", "LoggerMiddleware")))
	chain := middlewares.Chain(reqIdMiddleware, loggerMiddleware)

	indexHandler := handlers.NewIndexHandler(s.logger.With(slog.String("handler", "IndexHandler")))
	http.Handle("/", chain.Middleware(indexHandler))

	return http.ListenAndServe(":8787", nil)
}
