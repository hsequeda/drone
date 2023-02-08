package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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
	var dto LoadMedicationDTO
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
		if err == ErrNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	droneSerial := chi.URLParam(r, "serial")
	d, err := h.storage.Drone(droneSerial)
	if err != nil {
		if err == ErrNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	medDTOs := make([]MedicationDTO, len(d.Medications))
	for i, m := range d.Medications {
		// NOTE: use value convertion because the fields match for now.
		medDTOs[i] = MedicationDTO(m)
	}

	if err := json.NewEncoder(w).Encode(medDTOs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *DroneController) getAvailableDrones(w http.ResponseWriter, r *http.Request) {
	var availableDrones []AvailableDroneDTO
	for _, d := range h.storage.Drones() {
		if d.IsAvailable() {
			availableDrones = append(availableDrones, AvailableDroneDTO{
				Serial:          d.Serial,
				Model:           d.Model,
				WeightLimit:     d.WeightLimit,
				BatteryCapacity: d.BatteryCapacity,
				ConsumedWeight:  d.MedicationWeight(),
				State:           d.State,
			})
		}
	}

	if err := json.NewEncoder(w).Encode(availableDrones); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *DroneController) getDroneBatteryLevel(w http.ResponseWriter, r *http.Request) {
	droneSerial := chi.URLParam(r, "serial")
	d, err := h.storage.Drone(droneSerial)
	if err != nil {
		if err == ErrNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bc := DroneBatteryLevelDTO{BatteryLevel: d.BatteryCapacity}
	if err := json.NewEncoder(w).Encode(bc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
