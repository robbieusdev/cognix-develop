package oauth

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"time"
)

const (
	microsoftLoginURL = `https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=%s&scope=%s&response_type=code&redirect_uri=%s`
	microsoftToken    = `https://login.microsoftonline.com/organizations/oauth2/v2.0/token`
	/*
	   	client_id={client_id}&redirect_uri={redirect_uri}&client_secret={client_secret}
	   &code={code}&grant_type=authorization_code
	   	microsoftRefreshToken = ``

	*/
)

var microsoftScope = "offline_access Files.Read.All Sites.ReadWrite.All"

type (
	Config struct {
		Microsoft *MicrosoftConfig
		Google    *GoogleConfig
	}

	// MicrosoftConfig  declare configuration for Microsoft OAuth service
	MicrosoftConfig struct {
		ClientID     string `env:"MICROSOFT_CLIENT_ID,required"`
		ClientSecret string `env:"MICROSOFT_CLIENT_SECRET,required"`
		RedirectUL   string
	}
	microsoftExchangeCodeRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code,omitempty"`
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}
	tokenResponse struct {
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Microsoft implement OAuth authorization for microsoft`s services
	Microsoft struct {
		cfg        *MicrosoftConfig
		httpClient *resty.Client
	}
)

func (m *Microsoft) GetAuthURL(ctx context.Context, redirectUrl, state string) (string, error) {
	m.cfg.RedirectUL = fmt.Sprintf("%s/api/oauth/microsoft/callback", redirectUrl)
	return fmt.Sprintf(microsoftLoginURL, m.cfg.ClientID, microsoftScope, m.cfg.RedirectUL), nil
}

func (m *Microsoft) ExchangeCode(ctx context.Context, code string) (*IdentityResponse, error) {

	payload := map[string]string{
		"client_id":     m.cfg.ClientID,
		"client_secret": m.cfg.ClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  m.cfg.RedirectUL,
	}
	var response tokenResponse
	resp, err := m.httpClient.R().SetFormData(payload).
		Post(microsoftToken)
	if err != nil || resp.IsError() {
		return nil, utils.ErrorPermission.Newf("exchange code error %v : %v ", err, resp.Error())
	}
	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, err
	}
	return &IdentityResponse{
		Token: &oauth2.Token{
			AccessToken:  response.AccessToken,
			TokenType:    response.TokenType,
			RefreshToken: response.RefreshToken,
			Expiry:       time.Now().Add(time.Duration(response.ExpiresIn) * time.Second),
		},
	}, nil
}

func (m *Microsoft) RefreshToken(token *oauth2.Token) (*oauth2.Token, error) {
	payload := map[string]string{
		"client_id":     m.cfg.ClientID,
		"client_secret": m.cfg.ClientSecret,
		"refresh_token": token.RefreshToken,
		"grant_type":    "refresh_token",
		"redirect_uri":  m.cfg.RedirectUL,
	}
	var response tokenResponse
	resp, err := m.httpClient.R().SetFormData(payload).
		Post(microsoftToken)
	if err != nil || resp.IsError() {
		return nil, utils.ErrorPermission.Newf("exchange code error %v : %v ", err, resp.Error())
	}
	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken:  response.AccessToken,
		TokenType:    response.TokenType,
		RefreshToken: response.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(response.ExpiresIn) * time.Second),
	}, nil
}

func NewMicrosoft(cfg *MicrosoftConfig) Proxy {
	return &Microsoft{
		cfg:        cfg,
		httpClient: resty.New().SetTimeout(time.Minute),
	}
}
