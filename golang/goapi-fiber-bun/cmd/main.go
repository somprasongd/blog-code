package main

import (
	"fmt"
	"goapi/pkg/common/config"
	"goapi/pkg/common/database"
	"goapi/pkg/handlers"
	"goapi/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, closeDb, err := database.Init(cfg.Dsn)
	if err != nil {
		panic(err)
	}
	defer closeDb()

	app := fiber.New()

	// use middlewere
	app.Use(cors.New())
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(middlewares.Logger)

	handlers.RegisterRoutes(app, db)

	app.Listen(fmt.Sprintf(":%v", cfg.Port))
}
