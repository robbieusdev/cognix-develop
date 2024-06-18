package oauth

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
)

const (
	LoginState  = "login"
	SignUpState = "signUp"
	InviteState = "invite"

	ProviderGoogle    = "google"
	ProviderMicrosoft = "microsoft"
)

var Providers = map[string]bool{
	ProviderMicrosoft: true,
}

type (
	SignInConfig struct {
		State           string
		StateCookieName string
		URL             string
	}

	IdentityResponse struct {
		ID           string        `json:"id"`
		Email        string        `json:"email"`
		Name         string        `json:"name"`
		GivenName    string        `json:"given_name"`
		FamilyName   string        `json:"family_name"`
		AccessToken  string        `json:"access_token"`
		RefreshToken string        `json:"refresh_token"`
		Token        *oauth2.Token `json:"token,omitempty"`
	}

	Proxy interface {
		GetAuthURL(ctx context.Context, redirectUrl, state string) (string, error)
		ExchangeCode(ctx context.Context, code string) (*IdentityResponse, error)
		RefreshToken(token *oauth2.Token) (*oauth2.Token, error)
	}
)

func NewProvider(name string, cfg *Config) (Proxy, error) {
	switch name {
	case ProviderMicrosoft:
		return NewMicrosoft(cfg.Microsoft), nil
	}
	return nil, fmt.Errorf("unknown provider: %s", name)
}
