package models

type Trip struct {
	VehicleIdentifier string     `json:"vehicleIdentifier"`
	TripID            string     `json:"tripId"`
	TripGPS           []GPSPoint `json:"tripGPS"`
}

type GPSPoint struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Timestamp int64   `json:"ts"`
}

type TripsStore struct {
	Trips        map[string]Trip
	VehicleTrips map[string][]Trip
}

var Store = TripsStore{
	Trips:        make(map[string]Trip),
	VehicleTrips: make(map[string][]Trip),
}
