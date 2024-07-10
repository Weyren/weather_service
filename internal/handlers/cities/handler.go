package cities

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"weather_service/internal/database"
	"weather_service/pkg/utils"
)

const (
	citiesPath = "/api/cities"
	cityIDPath = "/api/cities/:id"
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
	router.HandlerFunc(http.MethodGet, citiesPath, h.GetAllCities)
}

func (h *Handler) GetAllCities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cities, err := h.repo.GetAllCities(r.Context())
	if err != nil {
		log.Println(err)
	}

	err = utils.WriteJSONIndented(w, cities)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error marshaling cities", http.StatusInternalServerError)
		return
	}

	log.Println("Get all cities", cities)

	w.WriteHeader(http.StatusOK)
}
