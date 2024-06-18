package handler

import (
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthHandler  handles authentication endpoints
type AuthHandler struct {
	oauthClient oauth.Proxy
	jwtService  security.JWTService
	authBL      bll.AuthBL
	//storage     storage.Storage
}

func NewAuthHandler(oauthClient oauth.Proxy,
	jwtService security.JWTService,
	authBL bll.AuthBL,
	//storage storage.Storage,

) *AuthHandler {
	return &AuthHandler{oauthClient: oauthClient,
		jwtService: jwtService,
		authBL:     authBL,
		//storage:    storage,
	}
}

func (h *AuthHandler) Mount(route *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := route.Group("/api/auth")
	handler.GET("/google/login", server.HandlerErrorFunc(h.SignIn))
	//handler.GET("/google/signup", server.HandlerErrorFunc(h.SignUp))
	handler.GET("/google/callback", server.HandlerErrorFunc(h.Callback))
	//handler.GET("/google/invite", server.HandlerErrorFunc(h.JoinToTenant))
	//adminHandler := route.Group("/api/auth")
	//adminHandler.Use(authMiddleware)
	//adminHandler.POST("/google/invite", server.HandlerErrorFuncAuth(h.Invite))
}

// SignIn login using google auth
// @Summary login using google auth
// @Description login using google auth
// @Tags Auth
// @ID auth_login
// @Param redirect_url query string false "redirect base url"
// @Produce  json
// @Success 200 {object} string
// @Router /auth/google/login [get]
func (h *AuthHandler) SignIn(c *gin.Context) error {
	var param parameters.LoginParam
	if err := c.ShouldBindQuery(&param); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong redirect url")
	}
	buf, err := json.Marshal(parameters.OAuthParam{Action: oauth.LoginState})
	if err != nil {
		return utils.Internal.Wrap(err, "can not marshal payload")
	}
	state := base64.URLEncoding.EncodeToString(buf)
	url, err := h.oauthClient.GetAuthURL(c.Request.Context(), param.RedirectURL, state)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, url)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) error {
	//buf, err := json.Marshal(parameters.OAuthParam{Action: oauth.LoginState})
	//if err != nil {
	//	return utils.Internal.Wrap(err, "can not marshal payload")
	//}
	//state := base64.URLEncoding.EncodeToString(buf)
	//url, err := h.oauthClient.RefreshToken(c.Request.Context(), state)
	//if err != nil {
	//	return err
	//}
	//c.Redirect(http.StatusFound, url)
	return nil
}

func (h *AuthHandler) Callback(c *gin.Context) error {
	code := c.Query(oauth.CodeNameGoogle)

	buf, err := base64.URLEncoding.DecodeString(c.Query(oauth.StateNameGoogle))
	if err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong state")
	}
	var state parameters.OAuthParam
	if err = json.Unmarshal(buf, &state); err != nil {
		return utils.Internal.Wrap(err, "can not unmarshal OAuth state")
	}

	response, err := h.oauthClient.ExchangeCode(c.Request.Context(), code)
	if err != nil {
		return err
	}
	var user *model.User
	switch state.Action {
	case oauth.LoginState:
		user, err = h.authBL.QuickLogin(c.Request.Context(), response)
	case oauth.SignUpState:
		user, err = h.authBL.SignUp(c.Request.Context(), response)
	//case oauth.InviteState:
	//	user, err = h.authBL.JoinToTenant(c.Request.Context(), &state, response)
	default:
		err = fmt.Errorf("unknown state %s ", state.Action)
	}
	if err != nil {
		return err
	}
	identity := &security.Identity{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		User: &model.User{ID: user.ID,
			TenantID: user.TenantID,
		},
	}
	token, err := h.jwtService.Create(identity)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, token)
}

//SignUp register new user and tenant using google auth
//@Summary register new user and tenant using google auth
//@Description register new user and tenant using google auth
//@Tags Auth
//@ID auth_signup
//@Produce  json
//@Success 200 {object} string
//@Router /auth/google/signup [get]
//func (h *AuthHandler) SignUp(c *gin.Context) error {
//	buf, err := json.Marshal(parameters.OAuthParam{Action: oauth.SignUpState})
//	if err != nil {
//		return utils.Internal.Wrap(err, "can not marshal payload")
//	}
//
//	state := base64.URLEncoding.EncodeToString(buf)
//	url, err := h.oauthClient.GetAuthURL(c.Request.Context(), state)
//	if err != nil {
//		return err
//	}
//	c.Redirect(http.StatusFound, url)
//	return nil
//}

// Invite create invitation for user
// @Summary create invitation for user
// @Description create invitation for user
// @Tags Auth
// @ID auth_invitation
// @Param params body parameters.InviteParam true "invitation  parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} string
// @Router /auth/google/invite [post]
//func (h *AuthHandler) Invite(c *gin.Context, identity *security.Identity) error {
//	var param parameters.InviteParam
//	if err := c.BindJSON(&param); err != nil {
//		return utils.ErrorBadRequest.Wrap(err, "can not parse payload")
//	}
//	if err := param.Validate(); err != nil {
//		return utils.ErrorBadRequest.Wrap(err, err.Error())
//	}
//
//	url, err := h.authBL.Invite(c.Request.Context(), identity, &param)
//	if err != nil {
//		return err
//	}
//	return server.JsonResult(c, http.StatusOK, url)
//}

// JoinToTenant join user to tenant using invitation link
// @Summary join user to tenant using invitation link
// @Description join user to tenant using invitation link
// @Tags Auth
// @ID auth_join_to_team
// @Produce  json
// @Success 200 {object} string
// @Router /auth/google/invite [get]
//func (h *AuthHandler) JoinToTenant(c *gin.Context) error {
//	//param := c.Query("state")
//	//
//	//key, err := base64.URLEncoding.DecodeString(param)
//	//if err != nil {
//	//	return utils.ErrorBadRequest.Wrap(err, "wrong state")
//	//}
//	////value, err := h.storage.Pull(string(key))
//	//
//	//state := base64.URLEncoding.EncodeToString("value")
//	//
//	//url, err := h.oauthClient.GetAuthURL(c.Request.Context(), state)
//	//if err != nil {
//	//	return err
//	//}
//	//c.Redirect(http.StatusFound, url)
//	return nil
//}
