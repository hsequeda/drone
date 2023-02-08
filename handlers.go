package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RegisterDroneDTO struct {
	Serial      string     `json:"serial"`
	Model       DroneModel `json:"model"`
	WeightLimit uint32     `json:"weight_limit"`
	Battery     uint8      `json:"battery"`
}

type MedicationDTO struct {
	Name   string
	Weight uint32
	Code   string
}

type DroneController struct {
	storage       *Storage
	maxUploadSize int64
}

func NewHttpServer(storage *Storage, maxUploadSize int64) *DroneController {
	return &DroneController{
		storage:       storage,
		maxUploadSize: maxUploadSize,
	}
}

func (h *DroneController) registerADrone(w http.ResponseWriter, r *http.Request) {
	dto := new(RegisterDroneDTO)
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	drone, err := NewDrone(dto.Serial, dto.Model, dto.WeightLimit, dto.Battery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.storage.SaveDrone(drone)
	json.NewEncoder(w).Encode("success")
}

func (h *DroneController) loadDrone(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.maxUploadSize)
	if err := r.ParseMultipartForm(h.maxUploadSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer func() { _ = file.Close() }()

	buff := make([]byte, 512)
	file.Read(buff)
	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		http.Error(w, "the provided file format is not allowed.", http.StatusBadRequest)
		return
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename, err := saveFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encryptedDto := r.PostFormValue("data")
	var dto MedicationDTO
	if err = json.Unmarshal([]byte(encryptedDto), &dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	meds, err := NewMedication(dto.Name, dto.Weight, dto.Code, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	droneSerial := chi.URLParam(r, "serial")
	drone, err := h.storage.Drone(droneSerial)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := drone.AddMedications(meds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.storage.SaveDrone(drone)
	json.NewEncoder(w).Encode("success")
}

func (h *DroneController) getDroneMedications(w http.ResponseWriter, r *http.Request) {
	panic("implement me!")
}

func (h *DroneController) getAvailableDrones(w http.ResponseWriter, r *http.Request) {
	panic("implement me!")
}

func (h *DroneController) getDroneBatteryLevel(w http.ResponseWriter, r *http.Request) {
	panic("implement me!")
}
