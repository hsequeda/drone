package storage

import (
	"context"
	"sync"
	"testing"

	"github.com/hsequeda/drone/drone"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// testInMemory is a shared storage used to test the parallelism in the integration test.
	testInMemory     *InMemory
	testInMemoryOnce sync.Once
)

var (
	// savedDroneSerial is the serial of a test Drone.
	savedDroneSerial = "45"
	// savedDroneWithMedicationSerial is the serial of a test Drone.
	savedDroneWithMedicationSerial = "50"
)

func TestInMemorySaveDrone(t *testing.T) {
	t.Parallel()
	s := initializeTestInMemory(t)
	d := drone.Drone{
		Serial:          "1",
		Model:           drone.Lightweight,
		WeightLimit:     300,
		BatteryCapacity: 100,
		State:           drone.Idle,
	}

	err := s.SaveDrone(context.Background(), d)
	require.NoError(t, err)
	savedDrone, _ := s.droneBySerial.Load(d.Serial)
	assert.Equal(t, d, savedDrone)
}

func TestInMemoryAddMedicationDrone(t *testing.T) {
	t.Parallel()
	s := initializeTestInMemory(t)
	val, _ := s.droneBySerial.Load(savedDroneWithMedicationSerial)
	d := val.(drone.Drone)
	// add medication
	newMedication := drone.Medication{
		Name:   "Aspirin",
		Weight: 50,
		Code:   "A01",
		Image:  "other_path",
	}
	d.Medications = append(d.Medications, newMedication)
	// initialize Drone
	err := s.SaveDrone(context.Background(), d)
	require.NoError(t, err)
	sd, _ := s.droneBySerial.Load(d.Serial)
	savedDrone := sd.(drone.Drone)
	assert.Equal(t, d, savedDrone)
	assert.Len(t, savedDrone.Medications, len(d.Medications))
}

func TestInMemoryGetDrone(t *testing.T) {
	t.Parallel()
	s := initializeTestInMemory(t)
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
			_, err := s.Drone(context.Background(), tc.serial)
			if tc.notFound {
				require.Error(t, err)
				assert.ErrorIs(t, err, drone.ErrNotFound)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestInMemoryGetDrones(t *testing.T) {
	t.Parallel()
	s := initializeTestInMemory(t)
	drones, err := s.Drones(context.Background())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(drones), 2) // compares with preset number of drones (could be more)
}

func initializeTestInMemory(t *testing.T) *InMemory {
	t.Helper()
	testInMemoryOnce.Do(func() {
		// setup storage
		testInMemory = NewInMemory()
		// add preset data for test
		testInMemory.droneBySerial.Store(savedDroneSerial, drone.Drone{
			Serial:          savedDroneSerial,
			Model:           drone.Lightweight,
			WeightLimit:     300,
			BatteryCapacity: 100,
			State:           drone.Idle,
		})

		testInMemory.droneBySerial.Store(savedDroneWithMedicationSerial, drone.Drone{
			Serial:          savedDroneWithMedicationSerial,
			Model:           drone.Heavyweight,
			WeightLimit:     400,
			BatteryCapacity: 90,
			State:           drone.Loaded,
			Medications: []drone.Medication{
				{
					Name:   "Omeprazol-250ml",
					Weight: 120,
					Code:   "OM_101",
					Image:  "/path/to/file",
				},
			},
		})
	})

	return testInMemory
}
