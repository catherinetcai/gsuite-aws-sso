package oauth

import "context"

// Service ...
type Service interface {
	GetOAuthLoginURL() string
	Exchange(ctx context.Context, code string) (*IDToken, error)
}
