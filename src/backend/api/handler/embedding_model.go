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

type EmbeddingModelHandler struct {
	embeddingModelBL bll.EmbeddingModelBL
}

func NewEmbeddingModelHandler(embeddingModelBL bll.EmbeddingModelBL) *EmbeddingModelHandler {
	return &EmbeddingModelHandler{embeddingModelBL: embeddingModelBL}
}

func (h *EmbeddingModelHandler) Mount(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	handler := router.Group("/api/manage/embedding_models").Use(authMiddleware)
	handler.GET("/", server.HandlerErrorFuncAuth(h.GetAll))
	handler.GET("/:id", server.HandlerErrorFuncAuth(h.GetByID))
	handler.POST("/", server.HandlerErrorFuncAuth(h.Create))
	handler.PUT("/:id", server.HandlerErrorFuncAuth(h.Update))
	handler.POST("/:id/:action", server.HandlerErrorFuncAuth(h.Delete))

}

// GetAll return list of embedding models
// @Summary return list of embedding models
// @Description return list of embedding models
// @Tags EmbeddingModel
// @ID embedding_model_get_all
// @Param archived query bool false "true for include deleted embedding models"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.EmbeddingModel
// @Router /manage/embedding_models [get]
func (h *EmbeddingModelHandler) GetAll(c *gin.Context, identity *security.Identity) error {
	var param parameters.ArchivedParam
	if err := c.ShouldBindQuery(&param); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong parameters")
	}
	result, err := h.embeddingModelBL.GetAll(c.Request.Context(), identity.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, result)
}

// GetByID return embedding model by id
// @Summary return embedding model by id
// @Description return embedding model by id
// @Tags EmbeddingModel
// @ID embedding_model_get_by_id
// @Param id path int true "embedding model id"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.EmbeddingModel
// @Router /manage/embedding_models/{id} [get]
func (h *EmbeddingModelHandler) GetByID(c *gin.Context, identity *security.Identity) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ErrorBadRequest.New("id should be presented")
	}
	embeddingModel, err := h.embeddingModelBL.GetByID(c.Request.Context(), identity.User, id)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, embeddingModel)
}

// Create creates embedding models
// @Summary creates embedding models
// @Description creates embedding models
// @Tags EmbeddingModel
// @ID embedding_model_create
// @Param params body parameters.EmbeddingModelParam true "embedding model parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 201 {object} model.EmbeddingModel
// @Router /manage/embedding_models [post]
func (h *EmbeddingModelHandler) Create(c *gin.Context, identity *security.Identity) error {
	var param parameters.EmbeddingModelParam
	if err := c.ShouldBind(&param); err != nil {
		return utils.ErrorBadRequest.New("invalid params")
	}
	embeddingModel, err := h.embeddingModelBL.Create(c.Request.Context(), identity.User, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusCreated, embeddingModel)
}

// Update updates embedding model
// @Summary updates embedding model
// @Description updates embedding model
// @Tags EmbeddingModel
// @ID embedding_model_update
// @Param id path int true "embedding model id"
// @Param params body parameters.EmbeddingModelParam true "embedding model parameter"
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.EmbeddingModel
// @Router /manage/embedding_models/{id} [put]
func (h *EmbeddingModelHandler) Update(c *gin.Context, identity *security.Identity) error {
	var param parameters.EmbeddingModelParam
	if err := c.ShouldBind(&param); err != nil {
		return utils.ErrorBadRequest.New("invalid params")
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ErrorBadRequest.New("id should be presented")
	}
	embeddingModel, err := h.embeddingModelBL.Update(c.Request.Context(), identity.User, id, &param)
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, embeddingModel)
}

// Delete deletes or restore embedding model
// @Summary delete or restore embedding model
// @Description delete or restore embedding model
// @Tags EmbeddingModel
// @ID embedding_model_delete
// @Param id path int true "embedding model id"
// @Param action path string true "action : delete | restore "
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} model.EmbeddingModel
// @Router /manage/embedding_models/{id}/{action} [post]
func (h *EmbeddingModelHandler) Delete(c *gin.Context, identity *security.Identity) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return utils.ErrorBadRequest.New("id should be presented")
	}
	action := c.Param("action")
	if !(action == ActionRestore || action == ActionDelete) {
		return utils.ErrorBadRequest.Newf("invalid action: should be %s or %s", ActionRestore, ActionDelete)
	}
	var embedingModel *model.EmbeddingModel
	switch action {
	case ActionRestore:
		embedingModel, err = h.embeddingModelBL.Restore(c.Request.Context(), identity.User, id)
	case ActionDelete:
		embedingModel, err = h.embeddingModelBL.Delete(c.Request.Context(), identity.User, id)
	}
	if err != nil {
		return err
	}
	return server.JsonResult(c, http.StatusOK, embedingModel)
}
