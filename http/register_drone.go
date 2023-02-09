package http

import (
	"encoding/json"
	"net/http"

	"github.com/hsequeda/drone/drone"
)

// RegisterDroneDTO struct is the value passed in the body of POST /registerDrone.
type RegisterDroneDTO struct {
	Serial      string      `json:"serial"`
	Model       drone.Model `json:"model"`
	WeightLimit uint32      `json:"weight_limit"`
	Battery     uint8       `json:"battery"`
}

func (h *DroneController) RegisterADrone(w http.ResponseWriter, r *http.Request) {
	dto := new(RegisterDroneDTO)
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d, err := drone.NewDrone(dto.Serial, dto.Model, dto.WeightLimit, dto.Battery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.storage.SaveDrone(r.Context(), d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode("success")
}
