package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"tripservice/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var mockMapboxServer *httptest.Server

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	mockMapboxServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockResponse := `{
			"features": [{
				"text": "SS16 5NP"
			}]
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))

	// Set the mock function for tests
	getPostcode = func(point models.GPSPoint) (string, error) {
		// Simulate Mapbox reverse geocode API response
		resp, err := http.Get(mockMapboxServer.URL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		var mapboxResp map[string]interface{}
		err = json.Unmarshal(body, &mapboxResp)
		if err != nil {
			return "", err
		}

		features := mapboxResp["features"].([]interface{})
		if len(features) > 0 {
			return features[0].(map[string]interface{})["text"].(string), nil
		}
		return "", nil
	}

	code := m.Run()
	mockMapboxServer.Close()
	os.Exit(code)
}

func TestAddTrips(t *testing.T) {
	router := gin.Default()
	router.POST("/trips", AddTrips)

	trips := []models.Trip{
		{
			VehicleIdentifier: "1",
			TripID:            "A",
			TripGPS: []models.GPSPoint{
				{Latitude: 51.558902, Longitude: 0.453003, Timestamp: 1615186800000},
				{Latitude: 51.558934, Longitude: 0.452974, Timestamp: 1615186801000},
			},
		},
	}
	jsonData, _ := json.Marshal(trips)
	req, _ := http.NewRequest("POST", "/trips", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTripPostcodes(t *testing.T) {
	router := gin.Default()
	router.GET("/trips/:id/postcodes", GetTripPostcodes)

	// Add a mock trip to the store
	models.Store.Trips["A"] = models.Trip{
		VehicleIdentifier: "1",
		TripID:            "A",
		TripGPS: []models.GPSPoint{
			{Latitude: 51.558902, Longitude: 0.453003, Timestamp: 1615186800000},
			{Latitude: 51.558934, Longitude: 0.452974, Timestamp: 1615186801000},
		},
	}

	req, _ := http.NewRequest("GET", "/trips/A/postcodes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SS16 5NP")
	assert.Contains(t, w.Body.String(), "SS16 5NP")
}
func TestGetTripSpeeds(t *testing.T) {
	router := gin.Default()
	router.GET("/trips/:id/speeds", GetTripSpeeds)

	// Add a mock trip to the store
	models.Store.Trips["A"] = models.Trip{
		VehicleIdentifier: "1",
		TripID:            "A",
		TripGPS: []models.GPSPoint{
			{Latitude: 51.558902, Longitude: 0.453003, Timestamp: 1615186800000},
			{Latitude: 51.558934, Longitude: 0.452974, Timestamp: 1615186801000},
			{Latitude: 51.559000, Longitude: 0.453050, Timestamp: 1615186802000},
		},
	}

	req, _ := http.NewRequest("GET", "/trips/A/speeds", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]float64
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	speeds := response["speeds"]
	assert.Len(t, speeds, 2)
	assert.True(t, speeds[0] > 0)
	assert.True(t, speeds[1] > 0)
}

func TestGetVehicleTrips(t *testing.T) {
	router := gin.Default()
	router.GET("/vehicles/:id/trips", GetVehicleTrips)

	// Add mock trips to the store for a vehicle
	models.Store.VehicleTrips["1"] = []models.Trip{
		{
			VehicleIdentifier: "1",
			TripID:            "A",
			TripGPS: []models.GPSPoint{
				{Latitude: 51.558902, Longitude: 0.453003, Timestamp: 1615186800000},
				{Latitude: 51.558934, Longitude: 0.452974, Timestamp: 1615186801000},
			},
		},
		{
			VehicleIdentifier: "1",
			TripID:            "B",
			TripGPS: []models.GPSPoint{
				{Latitude: 51.558900, Longitude: 0.453000, Timestamp: 1615186802000},
				{Latitude: 51.558940, Longitude: 0.453010, Timestamp: 1615186803000},
			},
		},
	}

	req, _ := http.NewRequest("GET", "/vehicles/1/trips", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	trips := response["trips"]
	assert.Len(t, trips, 2)

	assert.Equal(t, "A", trips[0]["trip_id"])
	assert.Equal(t, "SS16 5NP", trips[0]["start_postcode"])
	assert.Equal(t, "SS16 5NP", trips[0]["end_postcode"])
	assert.True(t, trips[0]["average_speed"].(float64) > 0)

	assert.Equal(t, "B", trips[1]["trip_id"])
	assert.Equal(t, "SS16 5NP", trips[1]["start_postcode"])
	assert.Equal(t, "SS16 5NP", trips[1]["end_postcode"])
	assert.True(t, trips[1]["average_speed"].(float64) > 0)
}
