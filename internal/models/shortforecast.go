package models

import "time"

// ShortForecast represents a short forecast for 1 day
type ShortForecast struct {
	Country  string
	City     string
	AvgTemp  float64
	DateList []time.Time
}
