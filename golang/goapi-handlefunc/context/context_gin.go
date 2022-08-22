package context

import (
	"github.com/gin-gonic/gin"
)

type GinContext struct {
	*gin.Context
}

func NewGinContext(c *gin.Context) MyContext {
	return &GinContext{
		Context: c,
	}
}

func (c *GinContext) Bind(v interface{}) error {
	return c.Context.ShouldBindJSON(v)
}

func (c *GinContext) BindQuery(v interface{}) error {
	return c.Context.ShouldBindQuery(v)
}

func (c *GinContext) Query(key string) (string, bool) {
	return c.Context.GetQuery(key)
}

func (c *GinContext) DefaultQuery(key string, d string) string {
	return c.Context.DefaultQuery(key, d)
}

func (c *GinContext) Param(key string) string {
	return c.Context.Param(key)
}

func (c *GinContext) Header(key string) string {
	return c.Context.GetHeader(key)
}

func (c *GinContext) RequestId() string {
	return c.Header("x-trace-id")
}
func (c *GinContext) ResponseError(code int, err string) {
	c.ResponseJSON(code, map[string]string{
		"error": err,
	})
}

func (c *GinContext) ResponseJSON(code int, data interface{}) {
	c.Context.JSON(code, data)
}

func WrapGinHandler(h func(MyContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(NewGinContext(c))
	}
}
