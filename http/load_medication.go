package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hsequeda/drone/drone"
)

const (

	// formPicture is the `Form Key` for the Medication Picture.
	formPicture = "picture"
	// formData is the `Form Key` for the Medication Payload.
	formData = "data"
)

// LoadMedicationDTO struct is the value passed in the body of PUT /drone/{serial}.
type LoadMedicationDTO struct {
	Name   string `json:"name"`
	Weight uint32 `json:"weight"`
	Code   string `json:"code"`
}

func (h *DroneController) LoadDrone(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.maxUploadSize)
	if err := r.ParseMultipartForm(h.maxUploadSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile(formPicture)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer func() { _ = file.Close() }()

	contentTypeBuff := make([]byte, 512)
	_, err = file.Read(contentTypeBuff)
	if err != nil {
		http.Error(w, fmt.Sprintf("read content-type buffer: %v", err), http.StatusBadRequest)
		return
	}
	filetype := http.DetectContentType(contentTypeBuff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		http.Error(w, "the provided file format is not allowed.", http.StatusBadRequest)
		return
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename, err := h.saveFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encryptedDto := r.PostFormValue(formData)
	var dto LoadMedicationDTO
	if err = json.Unmarshal([]byte(encryptedDto), &dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	meds, err := drone.NewMedication(dto.Name, dto.Weight, dto.Code, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

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

	if err := d.AddMedications(meds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.storage.SaveDrone(r.Context(), d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode("success")
}
