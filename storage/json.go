package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hsequeda/drone/drone"
	"github.com/sdomino/scribble"
)

const (
	// droneCollection const is the key for the drone collection in scribble db.
	droneCollection = "drone"
)

type JSON struct {
	db *scribble.Driver
}

func NewJSON(db *scribble.Driver) *JSON {
	return &JSON{db: db}
}

// Drone implements drone.Storage
func (j *JSON) Drone(ctx context.Context, serial string) (drone.Drone, error) {
	var d drone.Drone
	if err := j.db.Read(droneCollection, serial, &d); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return drone.Drone{}, drone.ErrNotFound
		}

		return drone.Drone{}, fmt.Errorf("read drone by serial: %w", err)
	}

	return d, nil
}

// Drones implements drone.Storage
func (j *JSON) Drones(ctx context.Context) ([]drone.Drone, error) {
	resp, err := j.db.ReadAll(droneCollection)
	if err != nil {
		return nil, errors.New("fetch all drones")
	}

	drones := make([]drone.Drone, len(resp))
	for i, v := range resp {
		var d drone.Drone
		if err = json.Unmarshal(v, &d); err != nil {
			return nil, fmt.Errorf("decode drone: %w", err)
		}

		drones[i] = d
	}

	return drones, nil
}

// SaveDrone implements drone.Storage
func (j *JSON) SaveDrone(ctx context.Context, d drone.Drone) error {
	if err := j.db.Write(droneCollection, d.Serial, d); err != nil {
		return fmt.Errorf("save drone: %w", err)
	}

	return nil
}

var _ drone.Storage = (*JSON)(nil)
