package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hsequeda/drone"
	dronehttp "github.com/hsequeda/drone/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// e2eSuite is a help struct to orchestate the e2e test.
// Similar to testify/testsuite, but simpler ;).
type e2eSuite struct {
	container  *DroneContainer
	testServer *httptest.Server
}

func TestE2E(t *testing.T) {
	s := new(e2eSuite)
	s.startServer(t)
	t.Cleanup(func() { s.testServer.Close() })

	t.Run("TestRegisterDrone", s.TestRegisterADrone)
	t.Run("TestAddMedication", s.TestAddMedication)
	t.Run("TestGetDroneMedications", s.TestGetDroneMedications)
	t.Run("TestGetAvailableDrones", s.TestGetAvailableDrones)
}

func (s *e2eSuite) TestRegisterADrone(t *testing.T) {
	t.Parallel()
	b, err := json.Marshal(dronehttp.RegisterDroneDTO{
		Serial:      "1",
		Model:       drone.Lightweight,
		WeightLimit: 300,
		Battery:     30,
	})
	require.NoError(t, err)
	resp, err := http.Post(s.buildURL("/drone"), "application/json", bytes.NewBuffer(b))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func (s *e2eSuite) TestAddMedication(t *testing.T) {
	t.Parallel()
	// setup storage data
	s.container.Storage().SaveDrone(drone.Drone{
		Serial:          "100",
		Model:           drone.Lightweight,
		WeightLimit:     400,
		BatteryCapacity: 80,
		State:           drone.Idle,
	})

	b, err := json.Marshal(dronehttp.LoadMedicationDTO{
		Name:   "Omeprazol-250g",
		Weight: 250,
		Code:   "OM_250",
	})

	require.NoError(t, err)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err = writer.WriteField("data", string(b))
	require.NoError(t, err)
	mediaPart, err := writer.CreateFormFile("picture", "OM_250")
	require.NoError(t, err)
	mediaData, err := ioutil.ReadFile("../../test/test_image.png")
	require.NoError(t, err)
	_, err = io.Copy(mediaPart, bytes.NewReader(mediaData))
	require.NoError(t, err)

	require.NoError(t, writer.Close())

	req, err := http.NewRequest(http.MethodPut, s.buildURL("/drone/100"), bytes.NewReader(body.Bytes()))
	req.Header.Add("Content-Type", writer.FormDataContentType())
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func (s *e2eSuite) TestGetDroneMedications(t *testing.T) {
	t.Parallel()
	// setup storage data
	m1 := drone.Medication{Name: "Omeprazol-250g", Weight: 250, Code: "OM_250", Image: "image_path"}
	m2 := drone.Medication{Name: "Advil-250g", Weight: 500, Code: "AD_500", Image: "another_image_path"}
	s.container.Storage().SaveDrone(drone.Drone{
		Serial:          "102",
		Model:           drone.Lightweight,
		WeightLimit:     400,
		BatteryCapacity: 80,
		State:           drone.Idle,
		Medications:     []drone.Medication{m1, m2},
	})

	resp, err := http.Get(s.buildURL("/drone/102/medications"))
	require.NoError(t, err)
	var body []dronehttp.MedicationDTO
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	assert.Len(t, body, 2)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	s.assertMedication(t, m1, body[0])
	s.assertMedication(t, m2, body[1])
}

func (s *e2eSuite) TestGetAvailableDrones(t *testing.T) {
	t.Parallel()
	// setup storage data
	availableD1 := drone.Drone{
		Serial:          "111",
		Model:           drone.Cruiserweight,
		WeightLimit:     400,
		BatteryCapacity: 80,
		State:           drone.Loading,
	}
	s.container.Storage().SaveDrone(availableD1)
	availableD2 := drone.Drone{
		Serial:          "112",
		Model:           drone.Cruiserweight,
		WeightLimit:     200,
		BatteryCapacity: 50,
		State:           drone.Idle,
	}
	s.container.Storage().SaveDrone(availableD2)
	unavailableD1 := drone.Drone{ // low battery
		Serial:          "113",
		Model:           drone.Cruiserweight,
		WeightLimit:     200,
		BatteryCapacity: 10,
		State:           drone.Idle,
	}
	s.container.Storage().SaveDrone(unavailableD1)
	unavailableD2 := drone.Drone{ // invalid state
		Serial:          "114",
		Model:           drone.Cruiserweight,
		WeightLimit:     200,
		BatteryCapacity: 95,
		State:           drone.Delivered,
	}
	s.container.Storage().SaveDrone(unavailableD2)
	unavailableD3 := drone.Drone{ // WeightLimit reached
		Serial:          "115",
		Model:           drone.Cruiserweight,
		WeightLimit:     250,
		BatteryCapacity: 95,
		State:           drone.Loading,
		Medications:     []drone.Medication{{Weight: 250}},
	}
	s.container.Storage().SaveDrone(unavailableD3)

	resp, err := http.Get(s.buildURL("/drones"))
	require.NoError(t, err)
	require.NoError(t, err)
	var body []dronehttp.AvailableDroneDTO
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(body), 2) // could be more than 3 (because we're sharing the Storage with the rest of the test)
	for _, add := range body {
		drone, err := s.container.Storage().Drone(add.Serial)
		require.NoError(t, err)
		s.assertAvailableDrone(t, drone, add)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func (s *e2eSuite) TestGetDroneBatteryLevel(t *testing.T) {
	t.Parallel()
	// setup storage data
	d := drone.Drone{
		Serial:          "1010",
		Model:           drone.Cruiserweight,
		WeightLimit:     400,
		BatteryCapacity: 80,
		State:           drone.Loading,
	}
	s.container.Storage().SaveDrone(d)

	resp, err := http.Get(s.buildURL("/drone/1010/battery"))
	require.NoError(t, err)
	var body dronehttp.DroneBatteryLevelDTO
	err = json.NewDecoder(resp.Body).Decode(&body)
	require.NoError(t, err)
	assert.Equal(t, 80, body.BatteryLevel)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func (s *e2eSuite) buildURL(path string) string {
	return s.testServer.URL + "/api/v1" + path
}

func (s *e2eSuite) startServer(t *testing.T) {
	t.Helper()
	s.container = NewDroneContainer(
		&Configuration{
			DroneController: DroneControllerConfiguration{
				MaxUploadSize: 5 * (1024 * 1024),
				UploadDir:     "~/tmp",
			},
		})

	s.testServer = httptest.NewServer(s.container.Router())
}

func (s *e2eSuite) assertMedication(t *testing.T, expected drone.Medication, actual dronehttp.MedicationDTO) bool {
	t.Helper()
	return assert.Equal(t, expected.Name, actual.Name) &&
		assert.Equal(t, expected.Weight, actual.Weight) &&
		assert.Equal(t, expected.Code, actual.Code) &&
		assert.Equal(t, expected.Image, actual.Image)
}

func (s *e2eSuite) assertAvailableDrone(t *testing.T, expected drone.Drone, actual dronehttp.AvailableDroneDTO) bool {
	t.Helper()
	return assert.Equal(t, expected.Model, actual.Model) &&
		assert.Equal(t, expected.State, actual.State) &&
		assert.Equal(t, expected.WeightLimit, actual.WeightLimit) &&
		assert.Equal(t, expected.BatteryCapacity, actual.BatteryCapacity) &&
		assert.Equal(t, expected.MedicationWeight(), actual.ConsumedWeight)
}
