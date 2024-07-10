package models

// City struct for OpenWeatherMap API
type City struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}
