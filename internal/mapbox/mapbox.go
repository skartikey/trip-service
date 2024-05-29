package mapbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"tripservice/models"
)

type MboxResponse struct {
	Features []struct {
		Text string `json:"text"`
	} `json:"features"`
}

func GetPostcode(point models.GPSPoint) (string, error) {
	mapboxToken := os.Getenv("mapboxToken")
	url := fmt.Sprintf("https://api.mapbox.com/geocoding/v5/mapbox.places/%f,%f.json?types=postcode&limit=1&access_token=%s", point.Longitude, point.Latitude, mapboxToken)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result MboxResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Features) == 0 {
		return "", fmt.Errorf("no postcode found")
	}

	return result.Features[0].Text, nil
}
