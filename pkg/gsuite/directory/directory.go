package discovery

import (
	"context"
	"errors"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/directory"
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
