package server

import (
	"goapi-hax/pkg/common/errs"
	"goapi-hax/pkg/core/ports"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinContext struct {
	*gin.Context
}

func NewGinContext(c *gin.Context) ports.Context {
	return &GinContext{
		Context: c,
	}
}

func (c *GinContext) Bind(v interface{}) error {
	return c.Context.ShouldBindJSON(v)
}

func (c *GinContext) Query(key string) (string, bool) {
	return c.Context.GetQuery(key)
}

func (c *GinContext) Param(key string) string {
	return c.Context.Param(key)
}

func (c *GinContext) Error(err error) {
	switch e := err.(type) {
	case errs.AppError:
		c.JSON(e.Code, e)
	case error:
		appErr := errs.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
		c.JSON(appErr.Code, appErr)
	}
}

func (c *GinContext) JSON(code int, data interface{}) {
	c.Context.JSON(code, data)
}

func NewGinHandler(handler func(ports.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewGinContext(c))
	}
}
