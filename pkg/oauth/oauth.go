package oauth

import (
	"context"
	"net/http"
	"net/url"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/config"
	httphelper "github.com/catherinetcai/gsuite-aws-sso/pkg/http"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	oauth2 "golang.org/x/oauth2"
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
	opts := defaultOptions()

	for _, setOpt := range setOpts {
		setOpt(opts)
	}

	return &Client{
		client:   opts.Client,
		logger:   opts.Logger,
		cfg:      oauthConf(opts.Config),
		loginURL: generateOAuthLoginURL(opts.Config),
	}
}

// LoginHandler redirects a client to the Google OAuth login page
func (c *Client) LoginHandler(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, c.loginURL, http.StatusFound)
}

// CallbackHandler handles the OAuth callback
// https://developers.google.com/identity/protocols/OAuth2WebServer
func (c *Client) CallbackHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	tok, err := c.cfg.Exchange(context.Background(), vars["code"])
	if err != nil {
		c.logger.Error("error exchanging OAuth code", zap.Error(err))
		httphelper.JSONResponse(w, struct{}{}, http.StatusBadRequest)
		return
	}

	// Get the ID token from the token, which when the 2nd portion is base64 decoded
	// gives us back something resembling:
	idTokenStr := tok.Extra("id_token").(string)

	idToken, err := ParseIDToken(idTokenStr)
	if err != nil {
		c.logger.Error("error parsing id token", zap.Error(err))
		httphelper.JSONResponse(w, struct{}{}, http.StatusBadRequest)
		return
	}

	// httphelper.JSONResponse(w, tok, http.StatusOK)
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
