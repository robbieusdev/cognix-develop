package oauth

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"time"
)

const (
	StateNameGoogle   = "state"
	CodeNameGoogle    = "code"
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type GoogleConfig struct {
	GoogleClientID string `env:"GOOGLE_CLIENT_ID"`
	GoogleSecret   string `env:"GOOGLE_SECRET"`
}

type googleProvider struct {
	config     *oauth2.Config
	httpClient *resty.Client
}

// NewGoogleProvider create new implementation of google oAuth client
func NewGoogleProvider(cfg *GoogleConfig, redirectURL string) Proxy {
	return &googleProvider{
		httpClient: resty.New().SetTimeout(time.Minute),
		config: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  fmt.Sprintf("%s/google/callback", redirectURL),
			Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile"},
		},
	}
}

func (g *googleProvider) GetAuthURL(ctx context.Context, redirectURL, state string) (string, error) {
	g.config.RedirectURL = fmt.Sprintf("%s/google/callback", redirectURL)

	return g.config.AuthCodeURL(state,
		oauth2.ApprovalForce), nil
}

func (g *googleProvider) ExchangeCode(ctx context.Context, code string) (*IdentityResponse, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, utils.Internal.Wrapf(err, "code exchange wrong: %s", err.Error())
	}
	response, err := g.httpClient.NewRequest().Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil || response.IsError() {
		return nil, utils.Internal.Newf("failed getting user info: %v [%v]", err.Error(), response.Error())
	}

	contents := response.Body()

	var data IdentityResponse
	if err = json.Unmarshal(contents, &data); err != nil {
		return nil, utils.Internal.Wrapf(err, "can not marshal google response")
	}
	data.AccessToken = token.AccessToken
	data.RefreshToken = token.RefreshToken

	return &data, nil
}

func (g *googleProvider) RefreshToken(token *oauth2.Token) (*oauth2.Token, error) {
	return token, nil
}
