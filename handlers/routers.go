package handlers

import (
	"github.com/gorilla/mux"
	"tranee_service/internal/logging"
	"tranee_service/services"
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
	r.HandleFunc("/countries/{id}", h.getOneCountry).Methods("GET")
	r.HandleFunc("/countries", h.getAllCountries).Methods("GET")
	r.HandleFunc("/countries", h.createCountry).Methods("POST")
	r.HandleFunc("/countries/{id}", h.changeCountry).Methods("PUT")
	r.HandleFunc("/countries/{id}", h.deleteCountry).Methods("DELETE")
	r.HandleFunc("/load-images", h.loadImages).Methods("GET")

	r.HandleFunc("/users", h.createUser).Methods("POST")
	r.HandleFunc("/users", h.getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", h.getUserById).Methods("GET")
	r.HandleFunc("/users/{id}", h.changeUser).Methods("PUT")
	r.HandleFunc("/users/{id}", h.deleteUser).Methods("DELETE")

	return r
}
