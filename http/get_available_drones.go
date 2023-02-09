package http

import (
	"encoding/json"
	"net/http"

	"github.com/hsequeda/drone/drone"
)

// AvailableDroneDTO struct is used in the response of GET /drones
type AvailableDroneDTO struct {
	Serial          string      `json:"serial"`
	Model           drone.Model `json:"model"`
	WeightLimit     uint32      `json:"weight_limit"`
	BatteryCapacity uint8       `json:"battery_capacity"`
	ConsumedWeight  uint32      `json:"consumed_weight"`
	State           drone.State `json:"state"`
}

func (h *DroneController) GetAvailableDrones(w http.ResponseWriter, r *http.Request) {
	var availableDrones []AvailableDroneDTO
	drones, err := h.storage.Drones(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, d := range drones {
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
