package server

import (
	"github.com/catherinetcai/gsuite-aws-sso/pkg/directory"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/logging"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/oauth"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/role"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Options contain all server options
type Options struct {
	Router    *mux.Router
	Logger    *zap.Logger
	Port      int
	OAuth     oauth.Service
	Directory directory.Service
	Role      role.Service
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

// WithOAuth sets the OAuth service
func WithOAuth(oAuth oauth.Service) Option {
	return func(o *Options) {
		o.OAuth = oAuth
	}
}

// WithDirectory sets the Directory Service
func WithDirectory(dir directory.Service) Option {
	return func(o *Options) {
		o.Directory = dir
	}
}

// WithRole sets the Role Service
func WithRole(r role.Service) Option {
	return func(o *Options) {
		o.Role = r
	}
}

func defaultOptions() *Options {
	return &Options{
		Logger: logging.Logger(),
		Port:   3030,
	}
}
