package handler

import (
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// TenantHandler  handles authentication endpoints
type TenantHandler struct {
	tenantBL bll.TenantBL
}

func NewTenantHandler(TenantBL bll.TenantBL) *TenantHandler {
	return &TenantHandler{
		tenantBL: TenantBL,
	}
}

func (h *TenantHandler) Mount(route *gin.Engine, TenantMiddleware gin.HandlerFunc) {
	handler := route.Group("/api/tenant")
	handler.Use(TenantMiddleware)
	handler.GET("/users", server.HandlerErrorFuncAuth(h.GetUserList))
	handler.GET("/user_info", server.HandlerErrorFuncAuth(h.GetUserInfo))
	handler.POST("/users", server.HandlerErrorFuncAuth(h.AddUser))
	handler.PUT("/users/:id", server.HandlerErrorFuncAuth(h.EditUser))
}

// GetUserList return list of users
// @Summary return list of users
// @Description return list of users
// @Tags Tenant
// @ID tenant_get_users
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.User
// @Router /tenant/users [get]
func (h *TenantHandler) GetUserList(c *gin.Context, identity *security.Identity) error {
	users, err := h.tenantBL.GetUsers(c.Request.Context(), identity.User)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, users)
}

// AddUser add new user
// @Summary add new user
// @Description add new user
// @Tags Tenant
// @ID tenant_add_user
// @Param params body parameters.AddUserParam true "create user parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.User
// @Router /tenant/users [post]
func (h *TenantHandler) AddUser(c *gin.Context, identity *security.Identity) error {
	if !identity.User.HasRoles(model.RoleAdmin, model.RoleSuperAdmin) {
		return utils.ErrorPermission.New("do not have permission")
	}

	var param parameters.AddUserParam
	if err := c.ShouldBind(&param); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong parameter")
	}
	if err := param.Validate(); err != nil {
		return utils.ErrorBadRequest.New(err.Error())
	}
	user, err := h.tenantBL.AddUser(c.Request.Context(), identity.User, param.Email, param.Role)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, user)
}

// EditUser edit  user
// @Summary edit user
// @Description edit  user
// @Tags Tenant
// @ID tenant_edit_user
// @Param id path string true "user id"
// @Param params body parameters.EditUserParam true "edit user parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.User
// @Router /tenant/users/{id} [put]
func (h *TenantHandler) EditUser(c *gin.Context, identity *security.Identity) error {
	if !identity.User.HasRoles(model.RoleAdmin, model.RoleSuperAdmin) {
		return utils.ErrorPermission.New("do not have permission")
	}
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return utils.ErrorPermission.Wrap(err, "wrong id parameter")
	}
	var param parameters.EditUserParam
	if err := c.ShouldBind(&param); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong parameter")
	}
	if err := param.Validate(); err != nil {
		return utils.ErrorBadRequest.New(err.Error())
	}
	user, err := h.tenantBL.UpdateUser(c.Request.Context(), identity.User, userID, param.Role)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, user)
}

// GetUserInfo get user info
// @Summary  get user info
// @Description  get user info
// @Tags Tenant
// @ID tenant_get_user_info
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.User
// @Router /tenant/user_info [get]
func (h *TenantHandler) GetUserInfo(c *gin.Context, identity *security.Identity) error {
	return server.JsonResult(c, http.StatusOK, identity.User)
}
