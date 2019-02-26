package oauth

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// IDToken matches the fields from the ID token of the OAuth response
type IDToken struct {
	Iss      string `json:"iss"`
	Aud      string `json:"aud"`
	Sub      string `json:"sub"`
	Hd       string `json:"hd"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	Iat      int    `json:"iat"`
	Exp      int    `json:"exp"`
}

// ParseIDToken takes in an ID token as string
func ParseIDToken(idToken string) (*IDToken, error) {
	split := strings.Split(idToken, ".")

	payload, err := base64.RawURLEncoding.DecodeString(split[1])
	if err != nil {
		return nil, err
	}

	token := &IDToken{}
	if err := json.Unmarshal(payload, token); err != nil {
		return nil, err
	}
	return token, nil
}
