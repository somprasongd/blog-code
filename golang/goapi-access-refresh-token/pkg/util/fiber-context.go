package util

import (
	"goapi/pkg/common"

	"github.com/gofiber/fiber/v2"
)

type fiberContext struct {
	*fiber.Ctx
}

func newFiberContext(c *fiber.Ctx) common.HContext {
	return &fiberContext{
		Ctx: c,
	}
}

func (c *fiberContext) Method() string {
	return c.Ctx.Method()
}

func (c *fiberContext) Path() string {
	return c.Ctx.Path()
}

func (c *fiberContext) BodyParser(v interface{}) error {
	return c.Ctx.BodyParser(v)
}

func (c *fiberContext) QueryParser(v interface{}) error {
	return c.Ctx.QueryParser(v)
}

func (c *fiberContext) Query(key string) (string, bool) {
	q := c.Ctx.Query(key)
	return q, true
}

func (c *fiberContext) DefaultQuery(key string, d string) string {
	return c.Ctx.Query(key, d)
}

func (c *fiberContext) Param(key string) string {
	return c.Ctx.Params(key)
}

func (c *fiberContext) Header(key string) string {
	return c.Ctx.Get(key)
}

func (c *fiberContext) Authorization() string {
	return c.Header("Authorization")
}

func (c *fiberContext) RequestId() string {
	return c.GetRespHeader("X-Request-ID")
}

func (c *fiberContext) Locals(key string, value ...interface{}) interface{} {
	return c.Ctx.Locals(key, value...)
}

func (c *fiberContext) Next() error {
	return c.Ctx.Next()
}

func (c *fiberContext) SendStatus(code int) error {
	return c.Ctx.SendStatus(code)
}

func (c *fiberContext) SendJSON(code int, data interface{}) error {
	c.Ctx.Status(code)
	return c.Ctx.JSON(data)
}

func WrapFiberHandler(h common.HandleFunc) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return h(newFiberContext(c))
	}
}
