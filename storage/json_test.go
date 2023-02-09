package storage

import (
	"context"
	"os"
	"testing"

	"github.com/hsequeda/drone/drone"
	"github.com/sdomino/scribble"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// jsonSuite is a help struct to orchestate the e2e test.
// Similar to testify/testsuite, but simpler ;).
type jsonSuite struct {
	db           *scribble.Driver
	storage      *JSON
	presetDrones []drone.Drone
}

func TestJSON(t *testing.T) {
	s := new(jsonSuite)
	db, err := scribble.New("../test/test_json_data", nil)
	require.NoError(t, err)
	s.db = db
	s.storage = NewJSON(db)
	s.presetDrones = append(s.presetDrones, drone.Drone{
		Serial:          "101",
		Model:           drone.Heavyweight,
		WeightLimit:     439,
		BatteryCapacity: 15,
		State:           drone.Idle,
	},
		drone.Drone{
			Serial:          "102",
			Model:           drone.Cruiserweight,
			WeightLimit:     100,
			BatteryCapacity: 98,
			State:           drone.Delivered,
		},
	)

	err = s.db.Write(droneCollection, s.presetDrones[0].Serial, s.presetDrones[0])
	require.NoError(t, err)
	err = s.db.Write(droneCollection, s.presetDrones[1].Serial, s.presetDrones[1])
	require.NoError(t, err)

	t.Cleanup(func() { os.RemoveAll("../test/test_json_data") })

	t.Run("TestSaveDrone", s.TestSaveDrone)
	t.Run("TestGetDrone", s.TestGetDrone)
	t.Run("TestGetDrones", s.TestGetDrones)
}

func (s *jsonSuite) TestSaveDrone(t *testing.T) {
	t.Parallel()
	d := drone.Drone{
		Serial:          "1",
		Model:           drone.Lightweight,
		WeightLimit:     300,
		BatteryCapacity: 100,
		State:           drone.Loaded,
		Medications: []drone.Medication{
			{
				Name:   "Omeprazol-250ml",
				Weight: 120,
				Code:   "OM_101",
				Image:  "/path/to/file",
			},
		},
	}

	err := s.storage.SaveDrone(context.Background(), d)
	require.NoError(t, err)

	var expected drone.Drone
	err = s.db.Read(droneCollection, d.Serial, &expected)
	require.NoError(t, err)
	assert.Equal(t, expected, d)
}

func (s *jsonSuite) TestGetDrone(t *testing.T) {
	t.Parallel()

	err := s.db.Write(droneCollection, s.presetDrones[0].Serial, s.presetDrones[0])
	require.NoError(t, err)

	testCases := []struct {
		name     string
		notFound bool
		serial   string
		expected drone.Drone
	}{
		{
			name:     "OK: existent drone",
			serial:   s.presetDrones[0].Serial,
			expected: s.presetDrones[0],
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
			d, err := s.storage.Drone(context.Background(), tc.serial)
			if tc.notFound {
				require.Error(t, err)
				assert.ErrorIs(t, err, drone.ErrNotFound)
				return
			}

			assert.Equal(t, tc.expected, d)
			require.NoError(t, err)
		})
	}
}

func (s *jsonSuite) TestGetDrones(t *testing.T) {
	t.Parallel()
	drones, err := s.storage.Drones(context.Background())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(drones), 2) // compares with preset number of drones (could be more)
}
