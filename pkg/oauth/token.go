package oauth

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// IDToken matches the fields from the ID token of the OAuth
type IDToken struct {
	Iss      string    `json:"iss"`
	Aud      string    `json:"aud"`
	Sub      string    `json:"sub"`
	Hd       string    `json:"hd"`
	Email    string    `json:"email"`
	Verified bool      `json:"verified"`
	Iat      time.Time `json:"iat"`
	Exp      time.Time `json:"exp"`
}

// ParseIDToken takes in an ID token as string
func ParseIDToken(idToken string) (*IDToken, error) {
	split := strings.Split(idToken, ".")

	payload, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		return nil, err
	}

	token := &IDToken{}
	if err := json.Unmarshal(payload, token); err != nil {
		return nil, err
	}
	return token, nil
}
