package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// testStorage is a shared storage used to test the parallelism in the integration test.
	testStorage     *Storage
	testStorageOnce sync.Once
)

var (
	// savedDroneSerial is the serial of a test Drone.
	savedDroneSerial = "45"
	// savedDroneWithMedicationSerial is the serial of a test Drone.
	savedDroneWithMedicationSerial = "50"
)

func TestSaveDrone(t *testing.T) {
	t.Parallel()
	s := initializeTestStorage(t)
	drone := Drone{
		Serial:          "1",
		Model:           Lightweight,
		WeightLimit:     300,
		BatteryCapacity: 100,
		State:           Idle,
	}

	s.SaveDrone(drone)
	savedDrone, _ := s.droneBySerial.Load(drone.Serial)
	assert.Equal(t, drone, savedDrone)
}

func TestAddMedicationDrone(t *testing.T) {
	t.Parallel()
	s := initializeTestStorage(t)
	d, _ := s.droneBySerial.Load(savedDroneWithMedicationSerial)
	drone := d.(Drone)
	// add medication
	newMedication := Medication{
		Name:   "Aspirin",
		Weight: 50,
		Code:   "A01",
		Image:  "other_path",
	}
	drone.Medications = append(drone.Medications, newMedication)
	// initialize Drone
	s.SaveDrone(drone)
	sd, _ := s.droneBySerial.Load(drone.Serial)
	savedDrone := sd.(Drone)
	assert.Equal(t, drone, savedDrone)
	assert.Len(t, savedDrone.Medications, len(drone.Medications))
}

func TestGetDrone(t *testing.T) {
	t.Parallel()
	s := initializeTestStorage(t)
	testCases := []struct {
		name     string
		notFound bool
		serial   string
	}{
		{
			name:   "OK: existent drone",
			serial: savedDroneSerial,
		},
		{
			name:     "Err: Not Found",
			serial:   "qwerty",
			notFound: true,
		},
	}

	for _, v := range testCases {
		tc := v
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := s.Drone(tc.serial)
			if tc.notFound {
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrNotFound)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGetDrones(t *testing.T) {
	t.Parallel()
	s := initializeTestStorage(t)
	assert.GreaterOrEqual(t, len(s.Drones()), 2) // compares with preset number of drones (could be more)
}

func initializeTestStorage(t *testing.T) *Storage {
	t.Helper()
	testStorageOnce.Do(func() {
		// setup storage
		testStorage = NewStorage()
		// add preset data for test
		testStorage.droneBySerial.Store(savedDroneSerial, Drone{
			Serial:          savedDroneSerial,
			Model:           Lightweight,
			WeightLimit:     300,
			BatteryCapacity: 100,
			State:           Idle,
		})

		testStorage.droneBySerial.Store(savedDroneWithMedicationSerial, Drone{
			Serial:          savedDroneWithMedicationSerial,
			Model:           Heavyweight,
			WeightLimit:     400,
			BatteryCapacity: 90,
			State:           Loaded,
			Medications: []Medication{
				{
					Name:   "Omeprazol-250ml",
					Weight: 120,
					Code:   "OM_101",
					Image:  "/path/to/file",
				},
			},
		})
	})

	return testStorage
}
