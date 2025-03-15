package geohelper

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"

	"github.com/tuan1kdt/soa-ba-test/internal/adapter/config"
)

type Geo struct {
	apiKey string
}

func New(config *config.GEO) *Geo {
	return &Geo{
		apiKey: config.APIKey,
	}
}

// Get location of a city using a geocoding API
func (g *Geo) GetCityLocation(city string) (float64, float64, error) {
	apiURL := fmt.Sprintf("https://api.opencagedata.com/geocode/v1/json?q=%s&key=%s", url.QueryEscape(city), g.apiKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Results []struct {
			Geometry struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"geometry"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, err
	}

	if len(data.Results) == 0 {
		return 0, 0, fmt.Errorf("no results found for city: %s", city)
	}

	return data.Results[0].Geometry.Lat, data.Results[0].Geometry.Lng, nil
}

func (g *Geo) GetIPLocation(ip string) (float64, float64, error) {
	// Example API: ip-api.com
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, err
	}

	return data.Lat, data.Lon, nil
}

func (g *Geo) GetDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth radius in kilometers
	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)

	lat1 = lat1 * (math.Pi / 180)
	lat2 = lat2 * (math.Pi / 180)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
