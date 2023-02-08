package drone

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMedication(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr bool

		medicationName   string
		medicationWeight uint32
		medicationCode   string
		medicationImage  string

		expected Medication
	}{
		{
			name:             "OK: Simple name",
			medicationName:   "Omeprazol",
			medicationWeight: 200,
			medicationCode:   "OM_101",
			medicationImage:  "/path/to/the_image",
			expected: Medication{
				Name:   "Omeprazol",
				Weight: 200,
				Code:   "OM_101",
				Image:  "/path/to/the_image",
			},
		},
		{
			name:             "OK: Name with numbers",
			medicationName:   "Omeprazol-250ml",
			medicationWeight: 200,
			medicationCode:   "OM_102",
			medicationImage:  "/path/to/the_image",
			expected: Medication{
				Name:   "Omeprazol-250ml",
				Weight: 200,
				Code:   "OM_102",
				Image:  "/path/to/the_image",
			},
		},
		{
			name:             "Err: 'name doesn't match'",
			expectedErr:      true,
			medicationName:   "0M3PA)*7",
			medicationWeight: 100,
			medicationCode:   "OM_101",
			medicationImage:  "/path/to/the_image",
		},
		{
			name:             "Err: 'code doesn't match'",
			expectedErr:      true,
			medicationName:   "Omeprazol",
			medicationWeight: 100,
			medicationCode:   "om_101", // code in lowercase
			medicationImage:  "/path/to/the_image",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newDrone, err := NewMedication(tc.medicationName, tc.medicationWeight, tc.medicationCode, tc.medicationImage)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, newDrone)
		})
	}
}
