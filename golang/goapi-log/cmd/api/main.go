package main

import (
	"goapi/pkg/common/logger"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	app := fiber.New()

	app.Use(requestid.New())
	app.Use(logMiddleware)

	app.Get("/", func(c *fiber.Ctx) error {
		log := c.Locals("log").(logger.Interface)
		log.Info("log in handler")
		return c.SendString("Hello, World!")
	})

	logger.Default.Info("start on port 3000")

	app.Listen(":3000")
}

func logMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	appName := os.Getenv("APP_NAME")

	if len(appName) == 0 {
		appName = "goapi"
	}

	fileds := map[string]interface{}{
		"app":       appName,
		"domain":    c.Hostname(),
		"requestId": c.GetRespHeader("X-Request-ID"),
		"userAgent": c.Get("User-Agent"),
		"ip":        c.IP(),
		"method":    c.Method(),
		"traceId":   c.Get("X-B3-Traceid"),
		"spanId":    c.Get("X-B3-Spanid"),
		"uri":       c.Path(),
	}

	log := logger.New(logger.ToFields(fileds)...)

	c.Locals("log", log)

	err := c.Next()

	fileds["status"] = c.Response().StatusCode()
	fileds["latency"] = time.Since(start)

	logger.Default.Info("", logger.ToFields(fileds)...)

	// logger.New(logger.ToFields(fileds)...).Info("")

	return err
}
