package geocoding

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"weather_service/internal/models"
	_ "weather_service/pkg/utils"
)

// GetCoordinates gets city coordinates from OpenWeatherMap API
func GetCoordinates(ctx context.Context, cityName string, appid string) (*models.City, error) {
	baseURL := "http://api.openweathermap.org/geo/1.0/direct?"

	params := url.Values{}
	params.Add("q", cityName)
	params.Add("limit", "5")
	params.Add("appid", appid)

	requestURL, _ := url.ParseRequestURI(baseURL)
	requestURL.RawQuery = params.Encode()

	//Request to OpenWeatherMap API with city name and appid key
	resp, err := http.Get(requestURL.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	//proxy := utils.Proxy{
	//	IP:   "47.251.70.179",
	//	Port: "80",
	//}
	//data, err := utils.GetResponseWithProxy(requestURL.String(), proxy)

	body, err := io.ReadAll(resp.Body)
	//Unmarshal JSON response
	var cities []models.City

	err = json.Unmarshal(body, &cities)
	if err != nil {
		return nil, err
	}

	if len(cities) == 0 {
		return nil, fmt.Errorf("city not found")
	}

	return &cities[0], nil
}

// GetCitiesCoordinates gets cities coordinates from OpenWeatherMap API
func GetCitiesCoordinates(cities []string, appid string) ([]*models.City, error) {

	var citiesCoords []*models.City
	for city := range cities {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		c, err := GetCoordinates(ctx, cities[city], appid)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		citiesCoords = append(citiesCoords, c)
	}
	return citiesCoords, nil

}

// UnmarshalJSON unmarshals JSON into Forecast struct and parses time
func UnmarshalJSON(b []byte, f *models.Forecast) (err error) {
	if err := json.Unmarshal(b, f); err != nil {
		return err
	}
	log.Println("forecast unmarshalled")

	for i := range f.List {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", f.List[i].DtTxt)
		if err != nil {
			return err
		}
		f.List[i].DtTime = parsedTime
	}

	return nil
}

// GetCityWeather gets city weather from OpenWeatherMap API
func GetCityWeather(city *models.City, appid string) (*models.Forecast, error) {
	baseURL := "http://api.openweathermap.org/data/2.5/forecast?"

	params := url.Values{}
	params.Add("lat", strconv.FormatFloat(city.Latitude, 'f', -1, 64))
	params.Add("lon", strconv.FormatFloat(city.Longitude, 'f', -1, 64))
	//params.Add("lang", "ru") - плохо работает
	params.Add("appid", appid)

	u, _ := url.ParseRequestURI(baseURL)
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var forecast models.Forecast
	err = UnmarshalJSON(body, &forecast)
	if err != nil {
		return nil, err
	}
	return &forecast, nil
}
