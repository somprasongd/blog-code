package main

import (
	"goapi-hax/pkg/common/config"
	"goapi-hax/pkg/common/database"
	"goapi-hax/pkg/core/services"
	"goapi-hax/pkg/handlers"
	"goapi-hax/pkg/repositories"
	"goapi-hax/pkg/server"
	"os"
)

func init() {

}

func main() {
	// load config
	config.LoadConfig()
	// For Liveness Probe
	if config.Config.App.Env == "production" {
		_, err := os.Create("/tmp/live")
		if err != nil {
			panic(err)
		}
		defer os.Remove("/tmp/live")
	}
	// connect database
	database.ConnectDB()
	// change NewTodoRepositoryMock to NewTodoRepositoryDB
	todoRepo := repositories.NewTodoRepositoryDB(database.DB)
	todoServ := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoServ)

	s := server.NewServer(todoHandler)

	s.Initialize()
}
