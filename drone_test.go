package drone

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDrone(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr bool

		droneSerial  string
		droneModel   DroneModel
		droneWeight  uint32
		droneBattery uint8

		expected Drone
	}{
		{
			name:         "OK: Lightweight drone",
			expectedErr:  false,
			droneSerial:  "1",
			droneModel:   Lightweight,
			droneWeight:  100,
			droneBattery: 80,
			expected: Drone{
				Serial:          "1",
				Model:           Lightweight,
				WeightLimit:     100,
				BatteryCapacity: 80,
				State:           Idle,
			},
		},
		{
			name:         "OK: Heavyweight drone",
			expectedErr:  false,
			droneSerial:  "2",
			droneModel:   Heavyweight,
			droneWeight:  200,
			droneBattery: 50,
			expected: Drone{
				Serial:          "2",
				Model:           Heavyweight,
				WeightLimit:     200,
				BatteryCapacity: 50,
				State:           Idle,
			},
		},
		{
			name:         "Err: 'weight limit exceed 500g'",
			expectedErr:  true,
			droneSerial:  "1",
			droneModel:   Lightweight,
			droneWeight:  800,
			droneBattery: 80,
		},
		{
			name:         "Err 'battery capacity exceed 100%'",
			expectedErr:  true,
			droneSerial:  "1",
			droneModel:   Lightweight,
			droneWeight:  800,
			droneBattery: 80,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newDrone, err := NewDrone(tc.droneSerial, tc.droneModel, tc.droneWeight, tc.droneBattery)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, newDrone)
		})
	}
}

func TestAddMedicationToDrone(t *testing.T) {
	om250g := Medication{
		Name:   "Omeprazol-250g",
		Weight: 250,
		Code:   "OM_250",
		Image:  "1023123asf",
	}
	testCases := []struct {
		name             string
		expectedErr      bool
		droneMedications []Medication
		droneBattery     uint8
		droneState       DroneState
		newMedication    Medication
	}{
		{
			name:          "OK-Idle",
			droneBattery:  80,
			droneState:    Idle,
			newMedication: om250g,
		},
		{
			name:          "OK-Loading",
			droneBattery:  80,
			droneState:    Loading,
			newMedication: om250g,
		},
		{
			name:          "Err-InvalidState",
			expectedErr:   true,
			droneBattery:  80,
			droneState:    Delivered,
			newMedication: om250g,
		},
		{
			name:          "Err-LowBattery",
			droneBattery:  20,
			expectedErr:   true,
			droneState:    Loading,
			newMedication: om250g,
		},
		{
			name:             "Err-TooMuchWeight",
			droneBattery:     80,
			expectedErr:      true,
			droneState:       Loading,
			droneMedications: []Medication{om250g},
			newMedication:    om250g,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newDrone := Drone{
				Serial:          "12345",
				Model:           Cruiserweight,
				WeightLimit:     400,
				BatteryCapacity: tc.droneBattery,
				State:           tc.droneState,
				Medications:     tc.droneMedications,
			}
			err := newDrone.AddMedications(tc.newMedication)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
