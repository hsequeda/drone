package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	dronehttp "github.com/hsequeda/drone/http"
	"github.com/hsequeda/drone/storage"
)

type Configuration struct {
	HTTPServer      HTTPServerConfiguration
	DroneController DroneControllerConfiguration
}

type DroneControllerConfiguration struct {
	MaxUploadSize int64
	UploadDir     string
}

type HTTPServerConfiguration struct {
	Addr string
}

type DroneContainer struct {
	config *Configuration

	router          *chi.Mux
	v1router        *chi.Mux
	httpServer      *http.Server
	storage         *storage.Storage
	droneController *dronehttp.DroneController
}

func NewDroneContainer(config *Configuration) *DroneContainer {
	return &DroneContainer{
		config: config,
	}
}

func (c *DroneContainer) Storage() *storage.Storage {
	if c.storage == nil {
		c.storage = storage.NewStorage()
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
		filesDir := http.Dir(c.config.DroneController.UploadDir)
		c.router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
			rctx := chi.RouteContext(r.Context())
			pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
			fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
			fs.ServeHTTP(w, r)
		})

		c.router.Mount("/debug", middleware.Profiler())
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
		c.v1router.Route("/", func(r chi.Router) {
			r.Get("/drones", c.DroneController().GetAvailableDrones)
			r.Post("/drone", c.DroneController().RegisterADrone)
			r.Put("/drone/{serial}", c.DroneController().LoadDrone)
			r.Get("/drone/{serial}/battery", c.DroneController().GetDroneBatteryLevel)
			r.Get("/drone/{serial}/medications", c.DroneController().GetDroneMedications)
		})
	}

	return c.v1router
}

func (c *DroneContainer) DroneController() *dronehttp.DroneController {
	if c.droneController == nil {
		c.droneController = dronehttp.NewHttpServer(c.Storage(), c.config.DroneController.MaxUploadSize, c.config.DroneController.UploadDir)
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
