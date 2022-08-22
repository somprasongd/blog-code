package main

import (
	"context"
	"flag"
	"fmt"
	"goapi/pkg/common/config"
	"goapi/pkg/common/database"
	"goapi/pkg/common/logger"
	"goapi/pkg/handlers"
	"goapi/pkg/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Test struct {
	Name string
}

func main() {
	// load config
	config.LoadConfig()
	// เรียกก่อนเริ่มเปิด server เพราะถ้าเชื่อมต่อไม่ได้ให้จะได้ไม่ต้อง start server
	database.ConnectDB()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// เปลี่ยนตรงนี้
	r := mux.NewRouter()
	// define route
	handler := setupRouter(r)

	// starting server
	port := config.Config.App.Port
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%v", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler, // Pass our instance of gorilla/mux in.
	}
	logger.Info(fmt.Sprintf("Starting server at port %v\n", port))
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Info(err.Error())
		}
	}()

	// Create channel to listen for signals.
	signalChan := make(chan os.Signal, 1)
	// Accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until receive signal.
	sig := <-signalChan
	logger.Info(fmt.Sprintf("%s signal caught", sig))

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Add extra handling here to clean up resources, such as flushing logs and
	// closing any database or Redis connections.

	// Gracefully shutdown the server by waiting on existing requests (except websockets).
	if err := srv.Shutdown(ctx); err != nil {
		logger.Info(fmt.Sprintf("server shutdown failed: %+v", err))
	}
	logger.Info("server exited")
	os.Exit(0)
}

func setupRouter(r *mux.Router) http.Handler {
	todo := r.PathPrefix("/api/todos").Subrouter()
	todoHandler := handlers.NewTodoHandler(database.DB)
	todo.HandleFunc("", todoHandler.CreateTodo).Methods(http.MethodPost)
	todo.HandleFunc("", todoHandler.ListTodo).Methods(http.MethodGet)
	// สามารถใช้ร่วมกับ regx ได้
	todo.HandleFunc("/{id:[0-9]+}", todoHandler.GetTodo).Methods(http.MethodGet)
	todo.HandleFunc("/{id:[0-9]+}", todoHandler.UpdateTodoStatus).Methods(http.MethodPut)
	todo.HandleFunc("/{id:[0-9]+}", todoHandler.DeleteTodo).Methods(http.MethodDelete)

	r.Use(middleware.Logging)

	// Handling CORS Requests
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
	})

	handler := c.Handler(r)
	return handler
}
