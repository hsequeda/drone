package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	execute(NewDroneContainer(&Configuration{
		HTTPServer: HTTPServerConfiguration{
			Addr: ":4444",
		},
		DroneController: DroneControllerConfiguration{
			MaxUploadSize: 5 * (1024 * 1024),
		},
	}))
}

func execute(c *DroneContainer) {
	debugRoutes(c.Router())
	go func() {
		if err := c.HTTPServer().ListenAndServe(); err != nil {
			fmt.Printf("error: %s", err.Error())
		}
	}()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		cancel()
	}()

	if err := c.HTTPServer().Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
		return
	}

	log.Println("server exited properly")
}

func debugRoutes(router *chi.Mux) {
	println("\nRoutes defined in the server:")
	chi.Walk(router, func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})
}
