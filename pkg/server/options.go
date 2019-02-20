package server

import (
	"github.com/catherinetcai/gsuite-aws-sso/pkg/logging"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Options contain all server options
type Options struct {
	Router *mux.Router
	Logger *zap.Logger
	Port   int
}

// Option is a functional way of setting options for the server
type Option func(o *Options)

// WithRouter sets a router on the Options struct
func WithRouter(r *mux.Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}

// WithLogger sets a logger on the Options struct
func WithLogger(l *zap.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// WithPort sets the port for the server to boot up
func WithPort(p int) Option {
	return func(o *Options) {
		o.Port = p
	}
}

func defaultOptions() *Options {
	return &Options{
		Logger: logging.Logger(),
		Port:   3030,
	}
}
