package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Server runs a HTTP server with some helpers
type Server struct {
	router *mux.Router
	logger *zap.Logger
	port   int
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
		router: opts.Router,
		logger: opts.Logger,
		port:   opts.Port,
	}, nil
}

// Run starts up the server
func (s *Server) Run() error {
	s.logger.Info("starting server", zap.Int("port", s.port))
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      handlers.RecoveryHandler(handlers.RecoveryLogger(&wrappedLogger{s.logger}))(s.router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return srv.ListenAndServe()
}

// HACK: Our ghetto fabulous wrapper for being able to have the zap logger work with the RecoveryLogger
type wrappedLogger struct {
	*zap.Logger
}

func (w *wrappedLogger) Println(args ...interface{}) {
	w.Logger.Sugar().Info(args...)
}
