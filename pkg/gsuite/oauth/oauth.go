package oauth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/config"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/oauth"
	"go.uber.org/zap"
	oauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// TODO: Need a way to exchange the refresh token for another access token
// Client encapsulates all OAuth actions
type Client struct {
	client   *http.Client
	logger   *zap.Logger
	loginURL string
	cfg      *oauth2.Config
}

// NewClient creates a new OAuth client
func NewClient(setOpts ...Option) *Client {
	opts := DefaultOptions()

	for _, setOpt := range setOpts {
		setOpt(opts)
	}

	return &Client{
		client: opts.Client,
		logger: opts.Logger,
		cfg:    oauthConf(opts.Config),
		// TODO: I don't love the reference the global config object, this breaks
		// the flow of data
		loginURL: generateOAuthLoginURL(opts.Config),
	}
}

// GetOAuthLoginURL returns the OAuth login URL
func (c *Client) GetOAuthLoginURL() string {
	return c.loginURL
}

// Exchange exchanges a code for a wrapped ID token
func (c *Client) Exchange(ctx context.Context, code string) (*oauth.IDToken, error) {
	tok, err := c.cfg.Exchange(context.Background(), code)
	if err != nil {
		c.logger.Error("error exchanging OAuth code", zap.Error(err))
		return nil, err
	}

	// Get the ID token from the token, which when the 2nd portion is base64 decoded
	// gives us back something resembling:
	idTokenStr := tok.Extra("id_token").(string)

	idToken, err := oauth.ParseIDToken(idTokenStr)
	if err != nil {
		c.logger.Error("error parsing id token", zap.Error(err))
		return nil, err
	}

	return idToken, nil
}

// TokenSourceFromCredentials takes in a credentials file and returns a token source
func (c *Client) TokenSourceFromCredentials(ctx context.Context, credentials []byte) (oauth2.TokenSource, error) {
	creds, err := google.CredentialsFromJSON(ctx, credentials, "email")
	if err != nil {
		c.logger.Error("error getting credentials from JSON", zap.Error(err))
		return nil, err
	}

	return creds.TokenSource, nil
}

// IDToken takes in a token source, validates it, and returns an ID token
// TODO: For now, we're not using the production version - we're just hitting an endpoint to
// validate that the token is valid. To productionize this we should be hitting the JWKS URI
func (c *Client) IDToken(tokenSource oauth2.TokenSource) (*oauth.IDToken, error) {
	v := url.Values{}
	token, err := tokenSource.Token()
	if err != nil {
		c.logger.Error("error getting token from source", zap.Error(err))
		return nil, err
	}

	idToken := token.Extra("id_token").(string)

	v.Set("id_token", idToken)

	url := url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "oauth2/v3/tokeninfo",
		RawQuery: v.Encode(),
	}

	resp, err := http.Get(url.String())
	if err != nil {
		c.logger.Error("error validating token", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("invalid request validating token")
	}

	return oauth.ParseIDToken(idToken)
}

func oauthConf(cfg config.OAuth) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       cfg.Scopes,
		Endpoint: oauth2.Endpoint{
			TokenURL: cfg.TokenURL,
			AuthURL:  cfg.AuthURL,
		},
		RedirectURL: cfg.RedirectURL,
	}
}

/*
	Formatting of the OAuth URL:
	https://accounts.google.com/o/oauth2/v2/auth?
	scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive.metadata.readonly&
	access_type=offline&
	include_granted_scopes=true&
	state=state_parameter_passthrough_value&
	redirect_uri=http%3A%2F%2Foauth2.example.com%2Fcallback&
	response_type=code&
	client_id=client_id
*/
func generateOAuthLoginURL(cfg config.OAuth) string {
	v := url.Values{}
	v.Set("scope", "email")
	// Setting access type to offline allows us to get an access and refresh token
	v.Set("access_type", "offline")
	v.Set("include_granted_scopes", "true")
	// TODO: Need to generate a nonce for the state passthrough value
	// v.Set("state", "state_parameter_passthrough_value")
	v.Set("redirect_uri", cfg.RedirectURL)
	v.Set("client_id", cfg.ClientID)
	v.Set("response_type", "code")

	// TODO: This should be extracted out
	url := url.URL{
		Scheme:   "https",
		Host:     "accounts.google.com",
		Path:     "o/oauth2/v2/auth",
		RawQuery: v.Encode(),
	}
	return url.String()
}
