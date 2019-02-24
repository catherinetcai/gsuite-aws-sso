package oauth

import (
	"net/http"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/config"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/logging"
	"go.uber.org/zap"
)

// Options contains all OAuthClient options
type Options struct {
	Logger *zap.Logger
	Client *http.Client
	Config config.OAuth
}

// Option is functional way of setting options for the OAuthClient
type Option func(o *Options)

// WithLogger sets logger on the Options struct
func WithLogger(l *zap.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// WithConfig sets the OAuth configs on the Options struct
func WithConfig(c config.OAuth) Option {
	return func(o *Options) {
		o.Config = c
	}
}

// WithClient sets an HTTP client on the Option
func WithClient(c *http.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// DefaultOptions set defaults
func DefaultOptions() *Options {
	return &Options{
		Client: &http.Client{},
		Logger: logging.Logger(),
	}
}
