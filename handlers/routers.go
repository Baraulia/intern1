package handlers

import (
	"github.com/gorilla/mux"
	"tranee_service/services"
)

type Handler struct {
	service services.AppCountries
	logger  services.Logger
}

func NewHandler(service services.AppCountries, logger services.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/countries/{id}", h.getOneCountry).Methods("GET")
	r.HandleFunc("/countries", h.getAllCountries).Methods("GET")

	return r
}
