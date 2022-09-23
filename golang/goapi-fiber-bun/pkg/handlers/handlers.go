package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func RegisterRoutes(app *fiber.App, db *bun.DB) {
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	api := app.Group("/api") // /api
	v1 := api.Group("/v1")   // /api/v1

	todoH := todoHandler{db}
	todo := v1.Group("/todos") // /api/v1/todos
	todo.Post("/", todoH.CreateTodo)
	todo.Get("/", todoH.ListTodo)
	todo.Get("/:id", todoH.GetTodo)
	todo.Patch("/:id", todoH.UpdateTodoStatus)
	todo.Delete("/:id", todoH.DeleteTodo)
}
