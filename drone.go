package main

import "errors"

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
