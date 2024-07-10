package geocoding

import (
	"context"
	"log"
	"sync"
	"time"
	"weather_service/internal/database"
	"weather_service/internal/models"
)

// WeatherUpdater - struct for updating weather data from OpenWeatherMap API every {interval} seconds
type WeatherUpdater struct {
	apikey string
	repo   database.Repository
	ticker *time.Ticker
}

// NewWeatherUpdater - constructor for WeatherUpdater struct
func NewWeatherUpdater(apikey string, repo database.Repository, interval time.Duration) *WeatherUpdater {
	return &WeatherUpdater{
		apikey: apikey,
		repo:   repo,
		ticker: time.NewTicker(interval),
	}
}

// Start - starts weather updater in background with {interval}
func (w *WeatherUpdater) Start() {
	go func() {
		defer w.ticker.Stop()
		for {
			select {
			case <-w.ticker.C:
				w.UpdateWeather()
			}
		}
	}()
}

// UpdateWeather asynchronously updates weather data from OpenWeatherMap API
func (w *WeatherUpdater) UpdateWeather() {
	ctx := context.Background()
	cities, err := w.repo.GetAllCities(ctx)
	if err != nil {
		log.Println(err)
	}
	wg := sync.WaitGroup{}
	for _, city := range cities {
		wg.Add(1)
		go func(city models.City) {
			defer wg.Done()
			log.Println("Fetching weather for city:", city)

			forecast, err := GetCityWeather(&city, w.apikey)
			if err != nil {
				log.Println(err)
			}

			//map for date and list of weather info for that date
			dateForecastMap := make(map[string][]models.List)

			for i := range forecast.List {
				//converting Kelvin to Celsius
				forecast.List[i].Main.Temp -= 273
				forecast.List[i].Main.FeelsLike -= 273
				forecast.List[i].Main.TempMin -= 273
				forecast.List[i].Main.TempMax -= 273

				date := forecast.List[i].DtTime.Format("2006-01-02")
				//creating map for date and list of weather info for that date
				if _, ok := dateForecastMap[date]; !ok {
					dateForecastMap[date] = make([]models.List, 0)
				}
				//appending weather info for that date
				dateForecastMap[date] = append(dateForecastMap[date], forecast.List[i])
			}

			for _, fk := range dateForecastMap {
				//creating weather info for that date
				finalWI := models.WeatherInfo{}

				var flagHasTemp bool = false
				for _, wi := range fk {
					//finding weather info for 12:00:00 time and setting temp for that date
					if wi.DtTime.Format("15:04:05") == "12:00:00" {
						finalWI.Temp = wi.Main.Temp
						finalWI.Date = wi.DtTime
						flagHasTemp = true

					}
					//appending additional info
					finalWI.AdditionalInfo = append(finalWI.AdditionalInfo, wi)
				}
				//if there is no weather info for 12:00:00 time, setting temp for earliest date
				if flagHasTemp == false {
					finalWI.Temp = fk[0].Main.Temp
				}
				//setting city id and date
				finalWI.CityID = city.ID
				finalWI.Date = fk[0].DtTime

				//saving weather info
				err = w.repo.CreateForecast(ctx, &finalWI, city.ID)
				if err != nil {
					log.Println(err)
				}
			}
		}(city)
	}
	//wait for all goroutines
	wg.Wait()
	log.Println("Weather updated")
}

// Stop - stops weather updater
func (w *WeatherUpdater) Stop() {
	w.ticker.Stop()
}
