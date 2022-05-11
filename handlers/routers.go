package handlers

import (
	"github.com/gorilla/mux"
	"lesson_2/pkg/logging"
	"lesson_2/services"
)

type Handler struct {
	service *services.Service
	logger  logging.Logger
}

func NewHandler(service *services.Service, logger logging.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/countries/:id", h.getAllCountries).Methods("GET")
	r.HandleFunc("/countries", h.getOneCountry).Methods("GET")

	return r
}
