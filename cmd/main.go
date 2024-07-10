package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
	"weather_service/internal/config"
	_ "weather_service/internal/database"
	"weather_service/internal/database/postgres"
	"weather_service/internal/geocoding"
	"weather_service/internal/handlers/cities"
	"weather_service/internal/handlers/forecasts"
	"weather_service/pkg/client"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not load config: error", err)
	}
	log.Println("Config loaded successfully")

	citiesList := []string{
		"Москва", "Санкт-Петербург", "Берлин", "Париж", "Лондон",
		"Токио", "Вашингтон", "Нью-Йорк", "Прага", "Будапешт",
		"Вена", "Мадрид", "Минск", "Севилья", "Рим",
		"Лиссабон", "Киев", "Белгород", "Калининград", "Амстердам",
	}

	router := httprouter.New()
	log.Println("Router created successfully")

	//postgres db client
	pgclient, err := client.NewClient(context.Background(), *cfg)
	if err != nil {
		log.Panicln("Can not connect DB: error", err)
	}
	log.Println("DB client connected successfully")

	//postgres repository
	repo := postgres.NewPostgresRepository(pgclient)
	log.Println("Repository created successfully")

	//Saving cities coordinates to table

	coords, err := geocoding.GetCitiesCoordinates(citiesList, cfg.API.Key)
	if err != nil {
		log.Println("Can not get coordinates: error", err)
	}
	log.Println("Coordinates received successfully")

	//Filling up table
	for _, city := range coords {
		err = repo.CreateCity(context.Background(), city)
		if err != nil {
			log.Println("Can not create city: error", err)
		}
	}
	log.Println("Table citiesList created and filled up successfully")

	//Weather updater init
	updater := geocoding.NewWeatherUpdater(cfg.API.Key, repo, time.Duration(cfg.API.Interval)*time.Second)
	updater.Start()
	updater.UpdateWeather()
	log.Println("Weather updater started")

	router.ServeFiles("/static/*filepath", http.Dir("static"))

	citiesHandler := cities.NewHandler(repo)
	citiesHandler.Register(router)

	forecastsHandler := forecasts.NewHandler(repo)
	forecastsHandler.Register(router)
	start(router, cfg)

}

func start(router *httprouter.Router, cfg *config.Config) {
	var listener net.Listener

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listen tcp on port", cfg.Server.Port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router}

	err = server.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
