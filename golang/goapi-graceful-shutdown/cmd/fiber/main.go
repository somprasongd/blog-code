package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	serverShutdown := make(chan struct{})

	go gracefulShutdown(app, serverShutdown)

	// Run the server
	port := 8080
	log.Printf("Starting server at port %v\n", port)

	err := app.Listen(fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil && err != http.ErrServerClosed {
		panic(err.Error())
	}

	select {
	case <-serverShutdown:
		log.Println("Shutdown completed")
	case <-time.After(timeout):
		log.Println("Shutdown timeout")
	}

	// <-serverShutdown
	log.Println("Running cleanup tasks")
	// Your cleanup tasks go here
}

func gracefulShutdown(srv *fiber.App, serverShutdown chan struct{}) {
	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	s := <-sig
	log.Printf("Received %v signal...", s)

	err := srv.Shutdown()
	if err != nil {
		log.Fatalf("Server shutdown failed: %+v\n", err)
	}
	serverShutdown <- struct{}{}
}
