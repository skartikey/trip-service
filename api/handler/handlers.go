package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"tripservice/internal/mapbox"
	"tripservice/models"
)

// Dependency injection for mapbox functions
var getPostcode = mapbox.GetPostcode

func AddTrips(c *gin.Context) {
	var newTrips []models.Trip
	if err := c.ShouldBindJSON(&newTrips); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, trip := range newTrips {
		models.Store.Trips[trip.TripID] = trip
		models.Store.VehicleTrips[trip.VehicleIdentifier] = append(models.Store.VehicleTrips[trip.VehicleIdentifier], trip)
	}

	c.JSON(http.StatusOK, gin.H{"status": "trips added"})
}

func GetTripPostcodes(c *gin.Context) {
	id := c.Param("id")
	trip, exists := models.Store.Trips[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
		return
	}

	startPostcode, err := getPostcode(trip.TripGPS[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get start postcode"})
		return
	}
	endPostcode, err := getPostcode(trip.TripGPS[len(trip.TripGPS)-1])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get end postcode"})
		return
	}

	// Using a struct to preserve the order of keys
	type PostcodesResponse struct {
		StartPostcode string `json:"start_postcode"`
		EndPostcode   string `json:"end_postcode"`
	}

	response := PostcodesResponse{
		StartPostcode: startPostcode,
		EndPostcode:   endPostcode,
	}

	c.JSON(http.StatusOK, response)
}

func GetTripSpeeds(c *gin.Context) {
	id := c.Param("id")
	trip, exists := models.Store.Trips[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
		return
	}

	speeds := make([]float64, len(trip.TripGPS)-1)
	for i := 1; i < len(trip.TripGPS); i++ {
		speeds[i-1] = calculateSpeed(trip.TripGPS[i-1], trip.TripGPS[i])
	}

	c.JSON(http.StatusOK, gin.H{"speeds": speeds})
}

func GetVehicleTrips(c *gin.Context) {
	id := c.Param("id")
	trips, exists := models.Store.VehicleTrips[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
		return
	}

	tripSummaries := make([]map[string]interface{}, len(trips))
	for i, trip := range trips {
		startPostcode, _ := getPostcode(trip.TripGPS[0])
		endPostcode, _ := getPostcode(trip.TripGPS[len(trip.TripGPS)-1])
		averageSpeed := calculateAverageSpeed(trip.TripGPS)
		tripSummaries[i] = map[string]interface{}{
			"trip_id":        trip.TripID,
			"start_postcode": startPostcode,
			"end_postcode":   endPostcode,
			"average_speed":  averageSpeed,
		}
	}

	c.JSON(http.StatusOK, gin.H{"trips": tripSummaries})
}

func calculateSpeed(p1, p2 models.GPSPoint) float64 {
	distance := haversine(p1.Latitude, p1.Longitude, p2.Latitude, p2.Longitude)
	timeDiff := float64(p2.Timestamp-p1.Timestamp) / 3600000.0 // milliseconds to hours
	speed := distance / timeDiff
	roundedSpeed, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", speed), 64)
	return roundedSpeed
}

func calculateAverageSpeed(points []models.GPSPoint) float64 {
	totalSpeed := 0.0
	for i := 1; i < len(points); i++ {
		totalSpeed += calculateSpeed(points[i-1], points[i])
	}
	averageSpeed := totalSpeed / float64(len(points)-1)
	roundedAverageSpeed, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", averageSpeed), 64)
	return roundedAverageSpeed
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers
	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
