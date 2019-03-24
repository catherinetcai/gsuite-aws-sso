package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/directory"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/oauth"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/role"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Server runs a HTTP server with some helpers
type Server struct {
	router       *mux.Router
	logger       *zap.Logger
	port         int
	oAuthSvc     oauth.Service
	directorySvc directory.Service
	roleSvc      role.Service
}

// New returns a new instance of the server
func New(setOpts ...Option) (*Server, error) {
	opts := defaultOptions()
	for _, setOpt := range setOpts {
		setOpt(opts)
	}

	if opts.Router == nil {
		return nil, fmt.Errorf("error: router cannot be nil")
	}

	return &Server{
		router:       opts.Router,
		logger:       opts.Logger,
		port:         opts.Port,
		oAuthSvc:     opts.OAuth,
		directorySvc: opts.Directory,
		roleSvc:      opts.Role,
	}, nil
}

// RegisterRoutes registers multiple routes
func (s *Server) RegisterRoutes(rs ...*Route) {
	for _, r := range rs {
		s.RegisterRoute(r)
	}
}

// RegisterRoute registers a single route
func (s *Server) RegisterRoute(r *Route) {
	s.logger.Info("Registering route",
		zap.String("path", r.Path),
		zap.String("method", r.Method.String()))

	s.router.HandleFunc(r.Path, r.HandlerFunc).Methods(r.Method.String()).Queries(r.Queries...)
}

// Run starts up the server
func (s *Server) Run() error {
	s.RegisterRoutes(s.defaultRoutes()...)

	s.logger.Info("starting server", zap.Int("port", s.port))
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      handlers.RecoveryHandler(handlers.RecoveryLogger(&wrappedLogger{Logger: s.logger}))(s.router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return srv.ListenAndServe()
}

// TODO: Not really happy with how these default routes are declared
func (s *Server) defaultRoutes() []*Route {
	return []*Route{
		&Route{
			Path:        "/auth/login",
			HandlerFunc: s.LoginHandler,
			Method:      GET,
		},
		&Route{
			Path:        "/auth/callback",
			HandlerFunc: s.CallbackHandler,
			Method:      GET,
			Queries:     []string{"code", "{code}"},
		},
		&Route{
			Path:        "/credentials",
			HandlerFunc: s.CredentialHandler,
			Method:      POST,
		},
		&Route{
			Path:        "/health",
			HandlerFunc: s.HealthHandler,
			Method:      GET,
		},
	}
}

// HACK: Our ghetto fabulous wrapper for being able to have the zap logger work with the RecoveryLogger
type wrappedLogger struct {
	Logger *zap.Logger
}

func (w *wrappedLogger) Println(args ...interface{}) {
	w.Logger.Sugar().Info(args...)
}
