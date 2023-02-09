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
	d := drone.Drone{
		Serial:          "1",
		Model:           drone.Lightweight,
		WeightLimit:     300,
		BatteryCapacity: 100,
		State:           drone.Idle,
	}

	s.SaveDrone(context.Background(), d)
	savedDrone, _ := s.droneBySerial.Load(d.Serial)
	assert.Equal(t, d, savedDrone)
}

func TestAddMedicationDrone(t *testing.T) {
	t.Parallel()
	s := initializeTestStorage(t)
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
	s.SaveDrone(context.Background(), d)
	sd, _ := s.droneBySerial.Load(d.Serial)
	savedDrone := sd.(drone.Drone)
	assert.Equal(t, d, savedDrone)
	assert.Len(t, savedDrone.Medications, len(d.Medications))
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

func TestGetDrones(t *testing.T) {
	t.Parallel()
	s := initializeTestStorage(t)
	drones, err := s.Drones(context.Background())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(drones), 2) // compares with preset number of drones (could be more)
}

func initializeTestStorage(t *testing.T) *Storage {
	t.Helper()
	testStorageOnce.Do(func() {
		// setup storage
		testStorage = NewStorage()
		// add preset data for test
		testStorage.droneBySerial.Store(savedDroneSerial, drone.Drone{
			Serial:          savedDroneSerial,
			Model:           drone.Lightweight,
			WeightLimit:     300,
			BatteryCapacity: 100,
			State:           drone.Idle,
		})

		testStorage.droneBySerial.Store(savedDroneWithMedicationSerial, drone.Drone{
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

	return testStorage
}
