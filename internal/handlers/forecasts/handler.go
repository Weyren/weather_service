package forecasts

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"weather_service/internal/database"
	"weather_service/pkg/utils"
)

const (
	shortForecastPath        = "/api/cities/:id/forecasts/shortforecast/"
	forecastPathWithDate     = "/api/cities/:id/forecasts/fullforecast/:date/"
	forecastPathWithDateTime = "/api/cities/:id/forecasts/fullforecast/:date/:time/"
)

type Handler struct {
	repo database.Repository
}

func NewHandler(repo database.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, shortForecastPath, h.GetShortForecastByCityID)
	router.HandlerFunc(http.MethodGet, forecastPathWithDate, h.GetForecastByCityIDandDate)
	router.HandlerFunc(http.MethodGet, forecastPathWithDateTime, h.GetForecastByCityIDandDateTime)
}

func (h Handler) GetShortForecastByCityID(w http.ResponseWriter, r *http.Request) {
	{
		w.Header().Set("Content-Type", "application/json")
		//w.Write([]byte("Forecast:\n"))

		params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)

		cityIDString := params.ByName("id")
		cityID, err := strconv.Atoi(cityIDString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		forecasts, err := h.repo.GetShortForecastByCityID(r.Context(), cityID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = utils.WriteJSONIndented(w, forecasts)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Println("Get short forecast", forecasts)
		w.WriteHeader(http.StatusOK)
	}
}

// GetForecastByCityIDandDate returns forecast for concrete date
func (h Handler) GetForecastByCityIDandDate(w http.ResponseWriter, r *http.Request) {
	{
		w.Header().Set("Content-Type", "application/json")
		params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)

		cityIDString := params.ByName("id")
		cityID, err := strconv.Atoi(cityIDString)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		date := params.ByName("date")

		forecasts, err := h.repo.GetForecastByCityIDandDate(r.Context(), cityID, date)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = utils.WriteJSONIndented(w, forecasts)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Println("Get forecast for concrete date", forecasts)
		w.WriteHeader(http.StatusOK)
	}
}

// GetForecastByCityIDandDateTime returns forecast for concrete date and time
func (h Handler) GetForecastByCityIDandDateTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)

	cityIDString := params.ByName("id")
	cityID, err := strconv.Atoi(cityIDString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	date := params.ByName("date")
	time := params.ByName("time")

	forecasts, err := h.repo.GetForecastByCityIDandDateTime(r.Context(), cityID, date, time)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = utils.WriteJSONIndented(w, forecasts)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Get forecast for concrete date and time", forecasts)
	w.WriteHeader(http.StatusOK)
}
