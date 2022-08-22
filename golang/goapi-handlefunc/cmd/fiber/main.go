package main

import (
	"goapi-handlefunc/context"
	"goapi-handlefunc/handler"

	"github.com/gofiber/fiber/v2"
)

const (
	BASE_URL = "/api/v1"
	PORT     = ":8080"
)

func main() {
	app := fiber.New()
	setRouter(app)

	app.Listen(PORT)
}

func setRouter(r *fiber.App) {
	h := handler.TodoHandler{}

	todos := r.Group(BASE_URL + "/todos")
	todos.Post("", context.WrapFiberHandler(h.CreateHandler))
	todos.Get("", context.WrapFiberHandler(h.ListHandler))
	todos.Get("/:id", context.WrapFiberHandler(h.GetHandler))
	todos.Patch("/:id", context.WrapFiberHandler(h.StatusUpdateHandler))
	todos.Delete("/:id", context.WrapFiberHandler(h.DeleteHandler))
}
