package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// droneSerialFromRequest extracts the drone serial from the path parameters.
func (h *DroneController) droneSerialFromRequest(r *http.Request) string {
	return chi.URLParam(r, "serial")
}
