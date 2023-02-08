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
}

func (s *e2eSuite) TestRegisterADrone(t *testing.T) {
	t.Parallel()
	b, err := json.Marshal(RegisterDroneDTO{
		Serial:      "1",
		Model:       Lightweight,
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
	s.container.Storage().SaveDrone(Drone{
		Serial:          "100",
		Model:           Lightweight,
		WeightLimit:     400,
		BatteryCapacity: 80,
		State:           Idle,
	})

	b, err := json.Marshal(MedicationDTO{
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
	mediaData, err := ioutil.ReadFile("./test/test_image.png")
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

func (s *e2eSuite) buildURL(path string) string {
	return s.testServer.URL + "/api/v1" + path
}

func (s *e2eSuite) startServer(t *testing.T) {
	t.Helper()
	s.container = NewDroneContainer(
		&Configuration{
			DroneController: DroneControllerConfiguration{
				MaxUploadSize: 5 * (1024 * 1024),
			},
		})

	s.testServer = httptest.NewServer(s.container.Router())
}
