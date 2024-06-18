package handler

import (
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
)

// OAuthHandler  provide oauth authentication for  third part services
type OAuthHandler struct {
	oauthConfig *oauth.Config
}

func NewOAuthHandler(oauthConfig *oauth.Config) *OAuthHandler {
	return &OAuthHandler{
		oauthConfig: oauthConfig,
	}
}

func (h *OAuthHandler) Mount(route *gin.Engine) {
	handler := route.Group("/api/oauth")
	handler.GET("/:provider/auth_url", server.HandlerErrorFunc(h.GetUrl))
	//handler.GET("/google/signup", server.HandlerErrorFunc(h.SignUp))
	handler.GET("/:provider/callback", server.HandlerErrorFunc(h.Callback))
	handler.POST("/:provider/refresh_token", server.HandlerErrorFunc(h.Refresh))
}

func (h *OAuthHandler) GetUrl(c *gin.Context) error {
	provider := c.Param("provider")
	var param parameters.LoginParam
	if err := c.ShouldBindQuery(&param); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong redirect url")
	}

	oauthClient, err := oauth.NewProvider(provider, h.oauthConfig)
	if err != nil {
		return utils.Internal.Wrap(err, "unknown provider")
	}
	url, err := oauthClient.GetAuthURL(c.Request.Context(), param.RedirectURL, "")
	if err != nil {
		return err
	}
	return server.StringResult(c, http.StatusOK, []byte(url))
}

func (h *OAuthHandler) Callback(c *gin.Context) error {
	provider := c.Param("provider")
	query := make(map[string]string)
	if err := c.BindQuery(&query); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong payload")
	}

	oauthClient, err := oauth.NewProvider(provider, h.oauthConfig)
	if err != nil {
		return utils.Internal.Wrap(err, "unknown provider")
	}
	result, err := oauthClient.ExchangeCode(c.Request.Context(), query["code"])
	if err != nil {
		return utils.ErrorPermission.New(err.Error())
	}
	return server.JsonResult(c, http.StatusOK, result)
}

func (h *OAuthHandler) Refresh(c *gin.Context) error {
	provider := c.Param("provider")
	var token oauth2.Token
	if err := c.BindJSON(&token); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong payload")
	}
	oauthClient, err := oauth.NewProvider(provider, h.oauthConfig)
	if err != nil {
		return utils.Internal.Wrap(err, "unknown provider")
	}

	result, err := oauthClient.RefreshToken(&token)
	if err != nil {
		return utils.ErrorPermission.New(err.Error())
	}
	_ = provider
	return server.JsonResult(c, http.StatusOK, result)
}
