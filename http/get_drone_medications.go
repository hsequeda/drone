package http

import (
	"encoding/json"
	"net/http"

	"github.com/hsequeda/drone/drone"
)

// MedicationDTO struct is used in the response of GET /drone/{serial}/medications
type MedicationDTO struct {
	Name   string `json:"name"`
	Weight uint32 `json:"weight"`
	Code   string `json:"code"`
	Image  string `json:"picture_path"`
}

func (h *DroneController) GetDroneMedications(w http.ResponseWriter, r *http.Request) {
	droneSerial := h.droneSerialFromRequest(r)
	d, err := h.storage.Drone(r.Context(), droneSerial)
	if err != nil {
		if err == drone.ErrNotFound {
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
