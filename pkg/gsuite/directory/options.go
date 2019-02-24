package discovery

import (
	"encoding/base64"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/logging"
	"go.uber.org/zap"
)

// Options ...
type Options struct {
	ClientID            string
	Logger              *zap.Logger
	ServiceAccountEmail string
	ServiceAccountPEM   []byte
}

// Option ...
type Option func(o *Options)

// WithServiceAccountBase64EncodedFile takes in the SA PEM file as a Base64 encoded string.
// Either this or the WithServiceAcountFile function must be used.
// If the base64 string is NOT valid, this WILL panic.
func WithServiceAccountBase64EncodedFile(file string) Option {
	return func(o *Options) {
		decoded, err := base64.StdEncoding.DecodeString(file)
		if err != nil {
			panic(err)
		}

		o.ServiceAccountPEM = decoded
	}
}

// WithServiceAccountFile takes in a file as bytes. Either this or
// the WithServiceAccountBase64EncodedFile will be used, depending on which one
// is run last.
func WithServiceAccountFile(f []byte) Option {
	return func(o *Options) {
		o.ServiceAccountPEM = f
	}
}

// WithServiceAccountEmail sets the email on the Option
func WithServiceAccountEmail(email string) Option {
	return func(o *Options) {
		o.ServiceAccountEmail = email
	}
}

// WithLogger sets the logger on the Option
func WithLogger(l *zap.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func defaultOptions() *Options {
	return &Options{
		Logger: logging.Logger(),
	}
}
