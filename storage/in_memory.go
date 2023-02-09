package storage

import (
	"context"
	"sync"

	"github.com/hsequeda/drone/drone"
)

// InMemory represents an 'In-Memory' storage for the service.
type InMemory struct {
	droneBySerial sync.Map
}

var _ drone.Storage = (*InMemory)(nil)

// NewInMemory initialize the Drone Storage.
func NewInMemory() *InMemory {
	return &InMemory{droneBySerial: sync.Map{}}
}

// Drone returns a Drone entity by its serial number.
// NOTE: Returns NotFound error if serial doesn't match.
func (s *InMemory) Drone(_ context.Context, serial string) (drone.Drone, error) {
	d, ok := s.droneBySerial.Load(serial)
	if !ok {
		return drone.Drone{}, drone.ErrNotFound
	}

	return d.(drone.Drone), nil
}

// Drone returns a list of all the Drone entities.
func (s *InMemory) Drones(_ context.Context) ([]drone.Drone, error) {
	droneArr := make([]drone.Drone, 0)
	s.droneBySerial.Range(func(_, d any) bool {
		droneArr = append(droneArr, d.(drone.Drone))
		return true
	})

	return droneArr, nil
}

// SaveDrone persists the current state of a Drone entity.
func (s *InMemory) SaveDrone(_ context.Context, drone drone.Drone) error {
	s.droneBySerial.Store(drone.Serial, drone)
	return nil
}
