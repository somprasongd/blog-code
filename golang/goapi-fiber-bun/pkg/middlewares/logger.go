package middlewares

import (
	"goapi/pkg/common/logger"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger(c *fiber.Ctx) error {
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

	log := logger.WithFields(fileds)

	c.Locals("log", log)

	err := c.Next()

	fileds["status"] = c.Response().StatusCode()
	fileds["latency"] = time.Since(start)

	logger.Info("access log", logger.ToFields(fileds)...)

	return err
}
