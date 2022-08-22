package server

import (
	"context"
	"fmt"
	"goapi-hax/pkg/common/config"
	"goapi-hax/pkg/common/logger"
	"goapi-hax/pkg/core/ports"
	"goapi-hax/pkg/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type server struct {
	todoHandler ports.TodoHandler
}

func NewServer(todoHandler ports.TodoHandler) *server {
	return &server{
		todoHandler: todoHandler,
	}
}

func (s server) Initialize() {
	r := mux.NewRouter()

	handler := s.setupRouter(r)

	port := config.Config.App.Port
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%v", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logger.Info(fmt.Sprintf("Starting server at port %v", port))
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
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

func (s server) setupRouter(r *mux.Router) http.Handler {
	r.HandleFunc("/healthz", handleHealthz).Methods(http.MethodGet)

	todo := r.PathPrefix("/api/todos").Subrouter()

	todo.HandleFunc("", NewHttpHandler(s.todoHandler.CreateTodo)).Methods(http.MethodPost)
	todo.HandleFunc("", NewHttpHandler(s.todoHandler.ListTodo)).Methods(http.MethodGet)
	// สามารถใช้ร่วมกับ regx ได้
	todo.HandleFunc("/{id:[0-9]+}", NewHttpHandler(s.todoHandler.GetTodo)).Methods(http.MethodGet)
	todo.HandleFunc("/{id:[0-9]+}", NewHttpHandler(s.todoHandler.UpdateTodo)).Methods(http.MethodPut)
	todo.HandleFunc("/{id:[0-9]+}", NewHttpHandler(s.todoHandler.DeleteTodo)).Methods(http.MethodDelete)

	// // Set up logging and panic recovery middleware.
	// amw := middleware.NewAuthenticationMiddleware()
	// todo.Use(amw.Middleware)
	r.Use(middleware.Logging)
	r.Use(middleware.PanicRecovery)

	// handle cors
	hc := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	// Insert the middleware
	handler := hc.Handler(r)
	return handler
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
