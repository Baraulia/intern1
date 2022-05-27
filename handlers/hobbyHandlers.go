package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"strconv"
	"tranee_service/models"
)

func (h *Handler) createHobby(w http.ResponseWriter, req *http.Request) {
	var input models.Hobby
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Errorf("Error while decoding request:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	result, err := govalidator.ValidateStruct(input)
	if !result {
		h.logger.Errorf("Incorrect data came from the request:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	hobbyId, err := h.service.AppHobbies.CreateHobby(&input)
	if err != nil {
		h.logger.Errorf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("id", strconv.Itoa(hobbyId))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getHobbies(w http.ResponseWriter, req *http.Request) {
	hobbies, err := h.service.AppHobbies.GetHobbies()
	if err != nil {
		h.logger.Warnf("server error: %s", err)
		http.Error(w, "server error", 500)
		return
	}
	output, err := json.Marshal(hobbies)
	if err != nil {
		h.logger.Errorf("getHobbies: error while marshaling list of hobbies: %s", err)
		http.Error(w, fmt.Sprintf("getHobbies: error while marshaling list of hobbies: %s", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getHobbies: error while writing response:%s", err)
		http.Error(w, fmt.Sprintf("getHobbies: error while writing response:%s", err), 500)
		return
	}
}
