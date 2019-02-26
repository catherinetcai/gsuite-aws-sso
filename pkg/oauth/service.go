package oauth

import (
	"context"

	"golang.org/x/oauth2"
)

// Service ...
type Service interface {
	GetOAuthLoginURL() string
	Exchange(ctx context.Context, code string) (*IDToken, error)
	TokenSourceFromCredentials(ctx context.Context, credentials []byte) (oauth2.TokenSource, error)
	IDToken(oauth2.TokenSource) (*IDToken, error)
}
