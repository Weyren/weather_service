package database

import (
	"context"
	"weather_service/internal/models"
)

type Repository interface {
	CreateCity(ctx context.Context, city *models.City) error
	GetAllCities(ctx context.Context) ([]models.City, error)
	CreateForecast(ctx context.Context, forecast *models.WeatherInfo, cityID int) error
	GetShortForecastByCityID(ctx context.Context, cityID int) (*models.ShortForecast, error)
	GetForecastByCityIDandDate(ctx context.Context, cityID int, datetime string) ([]models.WeatherInfo, error)
	GetForecastByCityIDandDateTime(ctx context.Context, cityID int, date string, time string) (*models.List, error)
}
