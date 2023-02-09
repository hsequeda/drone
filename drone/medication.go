package drone

import (
	"errors"
	"regexp"
)

// Validations for Medication struct fields.
var (
	nameValidation = regexp.MustCompile(`^([[:alpha:]]|[[:digit:]]|-|_)+$`)
	codeValidation = regexp.MustCompile(`^([[:upper:]]|[[:digit:]]|_)+$`)
)

// Medication defines the Medications to be carried by a drone.
type Medication struct {
	Name   string
	Weight uint32
	Code   string
	Image  string
}

// NewMedication builds a new instance of Medication.
func NewMedication(name string, weight uint32, code string, image string) (Medication, error) {
	if !nameValidation.MatchString(name) {
		return Medication{}, errors.New("name doesn't match")
	}

	if !codeValidation.MatchString(code) {
		return Medication{}, errors.New("code doesn't match")
	}

	return Medication{
		Name:   name,
		Weight: weight,
		Code:   code,
		Image:  image,
	}, nil
}
