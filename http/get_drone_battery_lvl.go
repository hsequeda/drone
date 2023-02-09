package http

import (
	"encoding/json"
	"net/http"

	"github.com/hsequeda/drone"
)

// DroneBatteryLevelDTO struct is used in the response of GET /drone/{serial}/battery_level
type DroneBatteryLevelDTO struct {
	BatteryLevel uint8 `json:"BatteryLevel"`
}

func (h *DroneController) GetDroneBatteryLevel(w http.ResponseWriter, r *http.Request) {
	droneSerial := h.droneSerialFromRequest(r)
	d, err := h.storage.Drone(droneSerial)
	if err != nil {
		if err == drone.ErrNotFound {
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
