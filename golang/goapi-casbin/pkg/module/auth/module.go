package auth

import (
	"goapi/pkg/app"
	"goapi/pkg/module/auth/core/ports"
	"goapi/pkg/module/auth/core/service"
	"goapi/pkg/module/auth/handler"
	"goapi/pkg/module/auth/repository"
	"goapi/pkg/util"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	BaseURL     string
	TokenSecret string
	Router      *fiber.App
	AuthService ports.AuthService
}

func Init(ctx *app.Context) {
	// สร้าง dependencies ทั้งหมด
	repo := repository.NewAuthRepositoryDB(ctx.DB.DB)
	svc := service.NewAuthService(ctx.Config, repo)

	cfg := RouteConfig{
		BaseURL:     ctx.Config.App.BaseUrl,
		TokenSecret: ctx.Config.Token.SecretKey,
		Router:      ctx.Router,
		AuthService: svc,
	}

	SetupRoutes(cfg)
}

func SetupRoutes(cfg RouteConfig) {
	h := handler.NewAuthHandler(cfg.AuthService)

	auth := cfg.Router.Group(cfg.BaseURL + "/auth")

	auth.Post("/register", util.WrapFiberHandler(h.Register))
	auth.Post("/login", util.WrapFiberHandler(h.Login))

	auth.Get("/profile", util.WrapFiberHandler(h.Profile))
	auth.Patch("/profile", util.WrapFiberHandler(h.UpdateProfile))
}
