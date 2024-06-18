package handler

import (
	"cognix.ch/api/v2/core/bll"
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ConnectorHandler struct {
	connectorBL bll.ConnectorBL
}

func NewCollectorHandler(connectorBL bll.ConnectorBL) *ConnectorHandler {
	return &ConnectorHandler{
		connectorBL: connectorBL,
	}
}
func (h *ConnectorHandler) Mount(route *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := route.Group("/api/manage/connector")
	handler.Use(authMiddleware)
	handler.GET("/source_types", server.HandlerErrorFuncAuth(h.GetSourceTypes))
	handler.GET("/", server.HandlerErrorFuncAuth(h.GetAll))
	handler.GET("/:id", server.HandlerErrorFuncAuth(h.GetById))
	handler.POST("/", server.HandlerErrorFuncAuth(h.Create))
	handler.PUT("/:id", server.HandlerErrorFuncAuth(h.Update))
	handler.POST("/:id/:action", server.HandlerErrorFuncAuth(h.Archive))
}

// GetAll return list of allowed connectors
// @Summary return list of allowed connectors
// @Description return list of allowed connectors
// @Tags Connectors
// @ID connectors_get_all
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.Connector
// @Router /manage/connector [get]
func (h *ConnectorHandler) GetAll(c *gin.Context, identity *security.Identity) error {
	connectors, err := h.connectorBL.GetAll(c.Request.Context(), identity.User)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, connectors)
}

// GetById return list of allowed connectors
// @Summary return list of allowed connectors
// @Description return list of allowed connectors
// @Tags Connectors
// @ID connectors_get_by_id
// @Produce  json
// @Param id path int true "connector id"
// @Security ApiKeyAuth
// @Success 200 {object} model.Connector
// @Router /manage/connector/{id} [get]
func (h *ConnectorHandler) GetById(c *gin.Context, identity *security.Identity) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ErrorBadRequest.New("id should be presented")
	}

	connectors, err := h.connectorBL.GetByID(c.Request.Context(), identity.User, id)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, connectors)
}

// Create creates connector
// @Summary creates connector
// @Description creates connector
// @Tags Connectors
// @ID connectors_create
// @Param params body parameters.CreateConnectorParam true "connector create parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 201 {object} model.Connector
// @Router /manage/connector/ [post]
func (h *ConnectorHandler) Create(c *gin.Context, identity *security.Identity) error {
	var param parameters.CreateConnectorParam
	if err := c.BindJSON(&param); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong payload")
	}
	if err := param.Validate(); err != nil {
		return utils.ErrorBadRequest.Wrap(err, err.Error())
	}
	connector, err := h.connectorBL.Create(c.Request.Context(), identity.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusCreated, connector)
}

// Update updates connector
// @Summary updates connector
// @Description updates connector
// @Tags Connectors
// @ID connectors_update
// @Param id path int true "connector id"
// @Param params body parameters.UpdateConnectorParam true "connector update parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.Connector
// @Router /manage/connector/{id} [put]
func (h *ConnectorHandler) Update(c *gin.Context, identity *security.Identity) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ErrorBadRequest.New("id should be presented")
	}
	var param parameters.UpdateConnectorParam
	if err = c.BindJSON(&param); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong payload")
	}
	if err = param.Validate(); err != nil {
		return utils.ErrorBadRequest.Wrapf(err, "validation error %s", err.Error())
	}
	connector, err := h.connectorBL.Update(c.Request.Context(), id, identity.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, connector)
}

// GetSourceTypes return list of source types
// @Summary return list of source types
// @Description return list of source types
// @Tags Connectors
// @ID connectors_get_source_types
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.SourceTypeDescription
// @Router /manage/connector/source_types [get]
func (h *ConnectorHandler) GetSourceTypes(c *gin.Context, identity *security.Identity) error {
	return server.JsonResult(c, http.StatusOK, model.SourceTypesList)
}

// Archive delete or restore connector
// @Summary delete or restore connector
// @Description delete or restore connector
// @Tags Connectors
// @ID Connectors_delete_restore
// @Param id path int true "Connectors id"
// @Param action path string true "action : delete | restore "
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.Connector
// @Router /manage/connector/{id}/{action} [post]
func (h *ConnectorHandler) Archive(c *gin.Context, identity *security.Identity) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ErrorBadRequest.New("id should be presented")
	}
	action := c.Param("action")
	if !(action == ActionRestore || action == ActionDelete) {
		return utils.ErrorBadRequest.Newf("invalid action: should be %s or %s", ActionRestore, ActionDelete)
	}
	connector, err := h.connectorBL.Archive(c.Request.Context(), identity.User, id, action == ActionRestore)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, connector)
}
