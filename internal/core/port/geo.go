package port

// GeoClient is an interface for interacting with a third-party geo service
type GeoClient interface {
	GetCityLocation(city string) (float64, float64, error)
	GetIPLocation(ip string) (float64, float64, error)
	GetDistance(lat1, lon1, lat2, lon2 float64) float64
}
