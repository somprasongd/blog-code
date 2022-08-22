package main

import (
	"goapi/pkg/app"
	"goapi/pkg/config"
	"goapi/pkg/module"

	"github.com/casbin/casbin/v2"
)

func main() {
	// Load config
	cfg := config.LoadConfig()
	// Load acl model and policy
	enforcer, err := casbin.NewEnforcer("config/acl_model.conf", "config/policy.csv")
	if err != nil {
		panic(err)
	}

	app := app.New(cfg)
	// Cleanup when server stopped
	defer app.Close()

	// For Liveness Probe
	app.CreateLivenessFile()

	// Initialize data sources
	app.InitDS()

	// Create router (mux/gin/fiber)
	app.InitRouter(enforcer)

	// Initialize module with dependency injection
	module.Init(app.Context)
	// Start server
	app.ServeHTTP()
}
