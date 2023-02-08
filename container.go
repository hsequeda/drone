package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Configuration struct {
	HTTPServer      HTTPServerConfiguration
	DroneController DroneControllerConfiguration
}

type DroneControllerConfiguration struct {
	MaxUploadSize int64
}

type HTTPServerConfiguration struct {
	Addr string
}

type DroneContainer struct {
	config *Configuration

	router          *chi.Mux
	v1router        *chi.Mux
	httpServer      *http.Server
	storage         *Storage
	droneController *DroneController
}

func NewDroneContainer(config *Configuration) *DroneContainer {
	return &DroneContainer{
		config: config,
	}
}

func (c *DroneContainer) Storage() *Storage {
	if c.storage == nil {
		c.storage = NewStorage()
	}

	return c.storage
}

func (c *DroneContainer) Router() *chi.Mux {
	if c.router == nil {
		c.router = chi.NewRouter()
		c.router.NotFound(c.router.NotFoundHandler())
		c.router.Route("/api", func(r chi.Router) {
			r.Mount("/v1", c.V1Router())
		})

		c.router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	}

	return c.router
}

func (c *DroneContainer) V1Router() *chi.Mux {
	if c.v1router == nil {
		c.v1router = chi.NewRouter()
		// Add middlewares
		c.v1router.Use(
			middleware.RequestID,
			middleware.Logger,
			middleware.Recoverer,
		)
		c.v1router.Route("/drone", func(r chi.Router) {
			r.Post("/", c.DroneController().registerADrone)
			r.Route("/{serial}", func(r chi.Router) {
				r.Put("/", c.DroneController().loadDrone)
			})
		})

	}

	return c.v1router
}

func (c *DroneContainer) DroneController() *DroneController {
	if c.droneController == nil {
		c.droneController = NewHttpServer(c.Storage(), c.config.DroneController.MaxUploadSize)
	}

	return c.droneController
}

func (c *DroneContainer) HTTPServer() *http.Server {
	if c.httpServer == nil {
		c.httpServer = &http.Server{
			Addr:    c.config.HTTPServer.Addr,
			Handler: c.Router(),
		}
	}
	return c.httpServer
}
