package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	httpAddr, ok := os.LookupEnv("HTTP_SERVER_ADDR")
	if !ok {
		log.Fatalf("HTTP_SERVER_ADDR is empty")
		return
	}

	uploadSizeStr, ok := os.LookupEnv("UPLOAD_SIZE")
	if !ok {
		log.Fatalf("UPLOAD_SIZE is empty")
		return
	}

	uploadSize, err := strconv.ParseInt(uploadSizeStr, 10, 64)
	if err != nil {
		log.Fatalf("UPLOAD_SIZE need to be integer")
		return
	}

	execute(NewDroneContainer(&Configuration{
		HTTPServer: HTTPServerConfiguration{
			Addr: httpAddr,
		},
		DroneController: DroneControllerConfiguration{
			MaxUploadSize: uploadSize * (1024 * 1024),
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
