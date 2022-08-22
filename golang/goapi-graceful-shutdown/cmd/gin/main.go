package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := gin.Default()
	// Add your routes as needed
	r.GET("/", func(ctx *gin.Context) {
		msg := "ok"
		ctx.Data(http.StatusOK, "text/plain", []byte(msg))
	})

	port := 8080
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%v", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	serverShutdown := make(chan struct{})

	go gracefulShutdown(srv, serverShutdown)

	log.Printf("Starting server at port %v\n", port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	select {
	case <-serverShutdown:
		log.Println("Shutdown completed")
	case <-time.After(timeout):
		log.Println("Shutdown timeout")
	}

	log.Println("Running cleanup tasks")
	// Your cleanup tasks go here
}

func gracefulShutdown(srv *http.Server, serverShutdown chan struct{}) {
	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	s := <-sig
	log.Printf("Received %v signal...", s)

	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("Server shutdown failed: %+v\n", err)
	}
	serverShutdown <- struct{}{}
}
