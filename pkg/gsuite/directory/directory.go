package directory

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/directory"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
)

const (
	awsSamlKey   = "AWS_SAML"
	defaultScope = "https://www.googleapis.com/auth/admin.directory.user"
)

var (
	ErrServiceAccountEmailNotSet = errors.New("service account email must be set")
	ErrServiceAccountFileNotSet  = errors.New("service account pem file must be set")
	ErrRoleNotSet                = errors.New("role not set on the user")
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

	// TODO: Add real scope here
	service, err := serviceClient(opts.ImpersonationEmail, opts.ServiceAccountPEM, []string{defaultScope, "https://www.googleapis.com/auth/admin.directory.group"})
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
	user, err := userSvc.Get(email).Projection("full").Do()
	if err != nil {
		c.logger.Error("error getting user", zap.Error(err))
		return nil, err
	}

	awsSamlInfo := &Attributes{}
	awsSamlInfoRaw, ok := user.CustomSchemas[awsSamlKey]
	if !ok {
		c.logger.Error("error attribute role info not found on user")
		return nil, ErrRoleNotSet
	}

	awsSamlInfoBytes, err := awsSamlInfoRaw.MarshalJSON()
	if err != nil {
		c.logger.Error("error marshalling raw JSON")
		return nil, err
	}

	err = json.Unmarshal(awsSamlInfoBytes, awsSamlInfo)
	if err != nil {
		c.logger.Error("error unmarshalling json", zap.Error(err))
		return nil, err
	}

	// TODO: Custom attributes to get mapped role arns
	return &directory.User{
		Email:        user.PrimaryEmail,
		CredentialID: getRole(awsSamlInfo),
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

func serviceClient(email string, credentials []byte, scopes []string) (*admin.Service, error) {
	config, err := google.JWTConfigFromJSON(credentials)
	if err != nil {
		return nil, err
	}

	config.Subject = email
	config.Scopes = scopes

	return admin.New(config.Client(context.Background()))
}

// TODO: Also get the session duration
func getRole(attributes *Attributes) string {
	roleInfo := strings.Split(attributes.IAMRole[0].Value, ",")
	return roleInfo[0]
}
