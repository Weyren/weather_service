package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"sort"
	"strings"
	"time"
	"weather_service/internal/database"
	"weather_service/internal/models"
	"weather_service/pkg/client"
)

// PostgresRepository implements Repository
type PostgresRepository struct {
	client client.Client
}

// NewPostgresRepository creates a new PostgresRepository returns interface
func NewPostgresRepository(client client.Client) database.Repository {
	return &PostgresRepository{
		client: client,
	}
}

// formatQuery formats query string for better logging
func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

// CreateCity creates a new city in database
func (r *PostgresRepository) CreateCity(ctx context.Context, city *models.City) error {
	q := `
	INSERT INTO cities
	    (city, country, lat, long) VALUES ($1, $2, $3, $4)
		RETURNING id
	
	`
	log.Println("SQL Query:", formatQuery(q), city.Name, city.Country, city.Latitude, city.Longitude)

	// insert new city
	if err := r.client.QueryRow(ctx, q, city.Name, city.Country, city.Latitude, city.Longitude).Scan(&city.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			log.Println(newErr)
		}
		return err
	}

	return nil
}

// GetAllCities returns all cities from database
func (r *PostgresRepository) GetAllCities(ctx context.Context) ([]models.City, error) {
	q := `SELECT id, city, country, lat, long FROM cities`
	// get all cities from database
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// create slice of cities
	cities := make([]models.City, 0)
	for rows.Next() {
		var city models.City
		// scan cities from database into slice of cities
		if err := rows.Scan(&city.ID, &city.Name, &city.Country, &city.Latitude, &city.Longitude); err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// sort cities by city name
	sort.Slice(cities, func(i, j int) bool {
		return cities[i].Name < cities[j].Name
	})

	return cities, nil
}

// CreateForecast creates a new forecast in database for concrete city
func (r *PostgresRepository) CreateForecast(ctx context.Context, forecast *models.WeatherInfo, cityID int) error {
	q := `INSERT INTO forecasts 
		(temp, date, additional_info, city_id) 
		VALUES ($1, $2, $3, $4)

		ON CONFLICT (city_id, date)
		DO UPDATE SET temp = excluded.temp, additional_info = excluded.additional_info
    `

	log.Println("SQL Query:", formatQuery(q), forecast.Temp, forecast.Date, forecast.AdditionalInfo, cityID)
	// insert new forecast for concrete city in database
	_, err := r.client.Exec(ctx, q, forecast.Temp, forecast.Date, forecast.AdditionalInfo, cityID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(pgErr)
		}
		return err
	}
	return nil
}

// GetShortForecastByCityID returns short forecast for concrete city
func (r *PostgresRepository) GetShortForecastByCityID(ctx context.Context, cityID int) (*models.ShortForecast, error) {
	q := `SELECT forecasts.temp, forecasts.date, cities.city, cities.country
		FROM forecasts
		JOIN cities ON cities.id = forecasts.city_id
		WHERE city_id = $1
		ORDER BY date
	`

	log.Println("SQL Query:", formatQuery(q), cityID)
	// get short forecast for 5 days for concrete city from database
	rows, err := r.client.Query(ctx, q, cityID)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}

	defer rows.Close()

	// scan short forecast from database
	var temp float64
	var sumTemp float64
	var count = 0

	var date time.Time
	var city, country string
	dateSlice := make([]time.Time, 0)

	for rows.Next() {
		if err := rows.Scan(&temp, &date, &city, &country); err != nil {
			log.Println("Scan error:", err)
			return nil, err
		}
		sumTemp += temp
		count++
		dateSlice = append(dateSlice, date)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error scanning row:", err)
		return nil, err
	}
	// calculate average temperature for 5 days
	avgtemp := sumTemp / float64(count)

	// create short forecast
	shortForecast := models.ShortForecast{
		City:     city,
		Country:  country,
		AvgTemp:  avgtemp,
		DateList: dateSlice,
	}

	return &shortForecast, nil

}

// GetForecastByCityIDandDate returns forecasts for concrete date
func (r *PostgresRepository) GetForecastByCityIDandDate(ctx context.Context, cityID int, date string) ([]models.WeatherInfo, error) {
	q := `
		SELECT id, temp, date, additional_info, city_id 
		FROM forecasts 
		WHERE city_id = $1
		AND date = $2
		ORDER BY date`

	log.Println("SQL Query:", formatQuery(q), cityID, date)
	// get forecasts for concrete date for concrete city from database
	rows, err := r.client.Query(ctx, q, cityID, date)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}

	defer rows.Close()

	// scan forecasts from database
	forecasts := make([]models.WeatherInfo, 0)
	for rows.Next() {
		var forecast models.WeatherInfo
		var additionalInfo json.RawMessage
		if err := rows.Scan(&forecast.ID, &forecast.Temp, &forecast.Date, &additionalInfo, &forecast.CityID); err != nil {
			log.Println("Scan error:", err)

			return nil, err
		}
		// unmarshal additional info includes forecasts for each 3 hours  from database
		if err := json.Unmarshal(additionalInfo, &forecast.AdditionalInfo); err != nil {
			log.Println("Unmarshal error:", err, " additionalInfo:", additionalInfo)
			return nil, err
		}
		forecasts = append(forecasts, forecast)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error scanning row:", err)
		return nil, err
	}

	return forecasts, nil

}

// GetForecastByCityIDandDateTime returns forecast for concrete date and time
func (r *PostgresRepository) GetForecastByCityIDandDateTime(ctx context.Context, cityID int, date, time string) (*models.List, error) {
	// get forecast for concrete date for concrete city from database
	forecasts, err := r.GetForecastByCityIDandDate(ctx, cityID, date)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// search for forecast for concrete time for concrete date
	for _, forecast := range forecasts {
		for _, info := range forecast.AdditionalInfo {
			if info.DtTime.Format("15:04:05") == time {
				return &info, nil
			}
		}
	}
	// if no forecast found
	err = errors.New("no forecast found")
	return nil, err
}
