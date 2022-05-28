package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
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
	r.HandleFunc("/countries/{id}", h.getOneCountry).Methods(http.MethodGet)
	r.HandleFunc("/countries", h.getAllCountries).Methods(http.MethodGet)
	r.HandleFunc("/countries", h.createCountry).Methods(http.MethodPost)
	r.HandleFunc("/countries/{id}", h.changeCountry).Methods(http.MethodPut)
	r.HandleFunc("/countries/{id}", h.deleteCountry).Methods(http.MethodDelete)
	r.HandleFunc("/load-images", h.loadImages).Methods(http.MethodGet)

	r.HandleFunc("/users", h.createUser).Methods(http.MethodPost)
	r.HandleFunc("/users", h.getUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", h.getUserById).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", h.changeUser).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", h.deleteUser).Methods(http.MethodDelete)
	r.HandleFunc("/users/{id}/hobbies", h.getHobbyByUserId).Methods(http.MethodGet)

	r.HandleFunc("/hobbies", h.createHobby).Methods(http.MethodPost)
	r.HandleFunc("/hobbies", h.getHobbies).Methods(http.MethodGet)

	return r
}
