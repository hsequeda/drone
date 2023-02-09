package drone

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("drone not found")

type Storage interface {
	// Drone returns a Drone entity by its serial number.
	// NOTE: Returns NotFound error if serial doesn't match.
	Drone(ctx context.Context, serial string) (Drone, error)
	// Drone returns a list of all the Drone entities.
	Drones(ctx context.Context) ([]Drone, error)
	// SaveDrone persists the current state of a Drone entity.
	SaveDrone(ctx context.Context, drone Drone) error
}
