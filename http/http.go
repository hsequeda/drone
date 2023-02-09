package http

import (
	"github.com/hsequeda/drone"
)

type DroneController struct {
	storage       *drone.Storage
	maxUploadSize int64
	uploadDir     string
}

func NewHttpServer(storage *drone.Storage, maxUploadSize int64, uploadDir string) *DroneController {
	return &DroneController{
		storage:       storage,
		maxUploadSize: maxUploadSize,
		uploadDir:     uploadDir,
	}
}
