package drone

import (
	"errors"
)

// DroneModel defines the different Models of drone.
type DroneModel int8

const (
	Lightweight DroneModel = iota + 1
	Middleweight
	Cruiserweight
	Heavyweight
)

// DroneState defines the different state availables in a drone.
type DroneState int8

const (
	Idle DroneState = iota + 1
	Loading
	Loaded
	Delivering
	Delivered
	Returning
)

// Drone defines the properties of a drone.
type Drone struct {
	Serial          string
	Model           DroneModel
	WeightLimit     uint32
	BatteryCapacity uint8
	State           DroneState
	Medications     []Medication
}

var (
	// ErrOverweight error occurs when the addition of a Medication exceed the `Weight Limit` of the drone.
	ErrOverweight = errors.New("unable to add medication: overweight")
	// ErrLowBattery error occurs when is tried 'to Load' a Drone with less than 25% of battery.
	ErrLowBattery = errors.New("unable to add medication: overweight")
	// ErrInvalidDroneState error occurs when is tried 'to Load' a Drone in a 'Loaded', 'Delivering', 'Delivered' or 'Returning' state.
	ErrInvalidDroneState = errors.New("invalid drone state")
)

// NewDrone builds a new IDLE drone instance.
func NewDrone(serial string, model DroneModel, weightLimit uint32, battery uint8) (Drone, error) {
	if weightLimit > 500 {
		return Drone{}, errors.New("weight limit exceed 500g")
	}

	if battery > 100 {
		return Drone{}, errors.New("battery capacity exceed 100%")
	}

	return Drone{
		Serial:          serial,
		Model:           model,
		WeightLimit:     weightLimit,
		BatteryCapacity: battery,
		State:           Idle,
	}, nil
}

// IsAvailable method returns if the current drone is available for load.
func (d *Drone) IsAvailable() bool {
	return (d.State == Idle ||
		d.State == Loading) &&
		d.BatteryCapacity > 25 &&
		d.MedicationWeight() < d.WeightLimit
}

// AddMedications method adds a new medication to the Drone if it doesn't
// exceed it WeightLimit.
func (d *Drone) AddMedications(m Medication) error {
	if d.State != Idle && d.State != Loading {
		return ErrInvalidDroneState
	}

	if d.BatteryCapacity < 25 {
		return ErrLowBattery
	}

	totalWeight := d.MedicationWeight() + m.Weight
	if d.WeightLimit < totalWeight {
		return ErrOverweight
	}

	d.Medications = append(d.Medications, m)
	return nil
}

// MedicationWeight method returns how many Weight is loading the drone.
func (d *Drone) MedicationWeight() uint32 {
	var w uint32
	for _, m := range d.Medications {
		w += m.Weight
	}
	return w
}
