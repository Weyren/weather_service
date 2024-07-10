package models

import "time"

// WeatherInfo struct represents weather info for saving to database
type WeatherInfo struct {
	ID             int       `json:"id"`
	Temp           float64   `json:"temp"`
	Date           time.Time `json:"date"`
	AdditionalInfo []List    `json:"additionalInfo"`
	CityID         int       `json:"city_id"`
}
