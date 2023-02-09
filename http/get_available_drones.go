package http

import (
	"encoding/json"
	"net/http"

	"github.com/hsequeda/drone"
)

// AvailableDroneDTO struct is used in the response of GET /drones
type AvailableDroneDTO struct {
	Serial          string           `json:"serial"`
	Model           drone.DroneModel `json:"model"`
	WeightLimit     uint32           `json:"weight_limit"`
	BatteryCapacity uint8            `json:"battery_capacity"`
	ConsumedWeight  uint32           `json:"consumed_weight"`
	State           drone.DroneState `json:"state"`
}

func (h *DroneController) GetAvailableDrones(w http.ResponseWriter, r *http.Request) {
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
