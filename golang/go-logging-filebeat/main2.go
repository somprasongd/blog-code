package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	fiberprometheus "github.com/somprasongd/fiber-prometheus"
)

func main() {
	fmt.Println("main")
	app := fiber.New()
	fmt.Println("NewMiddleware")
	p8sMiddleware := fiberprometheus.NewMiddleware("fiber", "http", "/metrics")
	fmt.Println("Register")
	p8sMiddleware.Register(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Go World!")
	})

	app.Get("/testurl", func(c *fiber.Ctx) error {
		return c.SendString("this is testurl.")
	})

	app.Get("hello/:name", func(c *fiber.Ctx) error {
		return c.SendString("Hello " + c.Params("name"))
	})
	fmt.Println("Listen")
	app.Listen(":7000")
}
