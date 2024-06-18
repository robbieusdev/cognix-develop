package server

import (
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context) error
type HandleFuncAuth func(c *gin.Context, identity *security.Identity) error
type JsonErrorResponse struct {
	Status        int    `json:"status,omitempty"`
	Error         string `json:"error,omitempty"`
	OriginalError string `json:"original_error,omitempty"`
}

type JsonResponse struct {
	Status int         `json:"status,omitempty"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func HandlerErrorFuncAuth(f HandleFuncAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity, err := GetContextIdentity(c)
		if err == nil {
			err = f(c, identity)
		}
		handleError(c, err)
	}
}

func HandlerErrorFunc(f HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handleError(c, f(c))
	}
}

func handleError(c *gin.Context, err error) {
	if err != nil {
		ew, ok := err.(utils.Errors)
		if !ok {
			ew.Original = err
			ew.Code = http.StatusInternalServerError
			ew.Message = err.Error()
		}
		errResp := JsonErrorResponse{
			Status: int(ew.Code),
			Error:  ew.Message,
		}
		if ew.Original != nil {
			errResp.OriginalError = ew.Original.Error()
		}
		zap.S().Errorf("[%s] %v", ew.Message, ew.Original)
		c.JSON(int(ew.Code), errResp)
	}
}

func JsonResult(c *gin.Context, status int, data interface{}) error {
	c.JSON(status, JsonResponse{
		Status: status,
		Error:  "",
		Data:   data,
	})
	return nil
}

func StringResult(c *gin.Context, status int, data []byte) error {
	c.Data(status, "", data)

	return nil
}
func BindJsonAndValidate(c *gin.Context, data interface{}) error {
	if err := c.BindJSON(data); err != nil {
		return utils.ErrorBadRequest.Wrap(err, "wrong payload")
	}
	if vl, ok := data.(validation.Validatable); ok {
		if err := vl.Validate(); err != nil {
			return utils.ErrorBadRequest.New(err.Error())
		}
	}
	return nil
}
