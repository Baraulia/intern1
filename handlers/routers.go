package handlers

import (
	"github.com/gorilla/mux"
	"tranee_service/internal/logging"
	"tranee_service/services"
)

type Handler struct {
	service services.AppCountries
	logger  logging.Logger
}

func NewHandler(service services.AppCountries, logger logging.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/countries/{id}", h.getOneCountry).Methods("GET")
	r.HandleFunc("/countries", h.getAllCountries).Methods("GET")
	r.HandleFunc("/countries", h.createCountry).Methods("POST")
	r.HandleFunc("/countries/{id}", h.changeCountry).Methods("PUT")
	r.HandleFunc("/countries/{id}", h.deleteCountry).Methods("DELETE")
	r.HandleFunc("/load-images", h.loadImages).Methods("GET")

	return r
}
