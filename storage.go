package main

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("drone not found")

// Storage represents an 'In-Memory' storage for the service.
type Storage struct {
	droneBySerial sync.Map
}

// NewStorage initialize the Drone Storage.
func NewStorage() *Storage {
	return &Storage{droneBySerial: sync.Map{}}
}

// Drone returns a Drone entity by its serial number.
// NOTE: Returns NotFound error if serial doesn't match.
func (s *Storage) Drone(serial string) (Drone, error) {
	drone, ok := s.droneBySerial.Load(serial)
	if !ok {
		return Drone{}, ErrNotFound
	}

	return drone.(Drone), nil
}

// Drone returns a list of all the Drone entities.
func (s *Storage) Drones() []Drone {
	droneArr := make([]Drone, 0)
	s.droneBySerial.Range(func(_, d any) bool {
		droneArr = append(droneArr, d.(Drone))
		return true
	})

	return droneArr
}

// SaveDrone persists the current state of a Drone entity.
func (s *Storage) SaveDrone(drone Drone) {
	s.droneBySerial.Store(drone.Serial, drone)
}
