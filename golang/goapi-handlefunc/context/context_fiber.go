package context

import (
	"github.com/gofiber/fiber/v2"
)

type FiberContext struct {
	*fiber.Ctx
}

func NewFiberContext(c *fiber.Ctx) MyContext {
	return &FiberContext{
		Ctx: c,
	}
}

func (c *FiberContext) Bind(v interface{}) error {
	return c.Ctx.BodyParser(v)
}

func (c *FiberContext) BindQuery(v interface{}) error {
	return c.Ctx.QueryParser(v)
}

func (c *FiberContext) Query(key string) (string, bool) {
	q := c.Ctx.Query(key)
	return q, true
}

func (c *FiberContext) DefaultQuery(key string, d string) string {
	return c.Ctx.Query(key, d)
}

func (c *FiberContext) Param(key string) string {
	return c.Ctx.Params(key)
}

func (c *FiberContext) Header(key string) string {
	return c.Ctx.Get(key)
}

func (c *FiberContext) RequestId() string {
	return c.Header("x-trace-id")
}

func (c *FiberContext) ResponseError(code int, err string) {
	c.ResponseJSON(code, map[string]string{
		"error": err,
	})
}

func (c *FiberContext) ResponseJSON(code int, data interface{}) {
	c.Ctx.SendStatus(code)
	c.Ctx.JSON(data)
}

func WrapFiberHandler(h func(MyContext)) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		h(NewFiberContext(c))
		return nil
	}
}
