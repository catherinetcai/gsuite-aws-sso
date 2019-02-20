package gsuite

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/directory"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/logging"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
)

var (
	ErrServiceAccountEmailNotSet = errors.New("service account email must be set")
	ErrServiceAccountFileNotSet  = errors.New("service account pem file must be set")
)

// Client is the GSuite client
type Client struct {
	service *admin.Service
	logger  *zap.Logger
}

// NewClient creates a new version of Client.
func NewClient(setOpts ...Option) (*Client, error) {
	opts := defaultOptions()

	for _, setOpt := range setOpts {
		setOpt(opts)
	}

	err := validateOpts(opts)
	if err != nil {
		return nil, err
	}

	service, err := serviceClient(opts.ServiceAccountEmail, opts.ServiceAccountPEM)
	if err != nil {
		return nil, err
	}

	return &Client{
		logger:  opts.Logger,
		service: service,
	}, nil
}

// GetUser find a user by email
func (c *Client) GetUser(email string) (*directory.User, error) {
	userSvc := admin.NewUsersService(c.service)
	user, err := userSvc.Get(email).Do()
	if err != nil {
		c.logger.Error("error getting user", zap.Error(err))
		return nil, err
	}

	// TODO: Custom attributes to get mapped role arns
	return &directory.User{
		Email: user.PrimaryEmail,
	}, nil
}

func validateOpts(opts *Options) error {
	if opts.ServiceAccountPEM == nil || len(opts.ServiceAccountPEM) == 0 {
		return ErrServiceAccountFileNotSet
	}

	if opts.ServiceAccountEmail == "" {
		return ErrServiceAccountEmailNotSet
	}
	return nil
}

func serviceClient(email string, credentials []byte) (*admin.Service, error) {
	config, err := google.JWTConfigFromJSON(credentials)
	if err != nil {
		return nil, err
	}

	config.Subject = email

	return admin.New(config.Client(context.Background()))
}

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

// WithServiceAccountEmail sets the email
func WithServiceAccountEmail(email string) Option {
	return func(o *Options) {
		o.ServiceAccountEmail = email
	}
}

func defaultOptions() *Options {
	return &Options{
		Logger: logging.Logger(),
	}
}
