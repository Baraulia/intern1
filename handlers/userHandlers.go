package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"strconv"
	"strings"
	"tranee_service/models"
)

func (h *Handler) createUser(w http.ResponseWriter, req *http.Request) {
	var input models.User
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
	userId, err := h.service.AppUsers.CreateUser(&input)
	if err != nil {
		h.logger.Errorf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("id", strconv.Itoa(userId))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getUsers(w http.ResponseWriter, req *http.Request) {
	var options models.Options
	if req.URL.Query().Get("page") != "" {
		paramPage, err := strconv.Atoi(req.URL.Query().Get("page"))
		if err != nil || paramPage < 0 {
			h.logger.Warnf("Invalid url request:%s", err)
			http.Error(w, "Invalid url request", 400)
			return
		}
		options.Page = uint64(paramPage)
	}
	if req.URL.Query().Get("limit") != "" {
		paramLimit, err := strconv.Atoi(req.URL.Query().Get("limit"))
		if err != nil || paramLimit < 0 {
			h.logger.Warnf("Invalid url request:%s", err)
			http.Error(w, "Invalid url request", 400)
			return
		}
		options.Limit = uint64(paramLimit)
	}
	users, pages, err := h.service.AppUsers.GetUsers(&options)
	if err != nil {
		h.logger.Warnf("server error: %s", err)
		http.Error(w, "server error", 500)
		return
	}
	output, err := json.Marshal(users)
	if err != nil {
		h.logger.Errorf("getUsers: error while marshaling list of users: %s", err)
		http.Error(w, fmt.Sprintf("getUsers: error while marshaling list of users: %s", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Pages", strconv.Itoa(pages))
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getUsers: error while writing response:%s", err)
		http.Error(w, fmt.Sprintf("getUsers: error while writing response:%s", err), 500)
		return
	}
}

func (h *Handler) getUserById(w http.ResponseWriter, req *http.Request) {
	paramId := strings.TrimPrefix(req.URL.Path, "/users/")
	userId, err := strconv.Atoi(paramId)
	if err != nil || userId <= 0 {
		h.logger.Warnf("Invalid request:%s", err)
		http.Error(w, "Invalid url request", 400)
		return
	}
	user, err := h.service.AppUsers.GetUserById(userId)
	if err != nil {
		if err.Error() == "such a user does not exist" {
			h.logger.Warnf("getUserById: such user does not exist")
			http.Error(w, "such user does not exist", 404)
			return
		}
		h.logger.Warnf("getUserById: server error: %s", err)
		http.Error(w, "server error", 500)
		return
	}
	output, err := json.Marshal(user)
	if err != nil {
		h.logger.Errorf("getUserById: error while marshaling one user: %s", err)
		http.Error(w, fmt.Sprintf("getUserById: error while marshaling one user: %s", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getUserById: error while writing response:%s", err)
		http.Error(w, fmt.Sprintf("getUserById: error while writing response:%s", err), 500)
		return
	}
}

func (h *Handler) changeUser(w http.ResponseWriter, req *http.Request) {
	var input models.User
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
	paramId := strings.TrimPrefix(req.URL.Path, "/users/")
	userId, err := strconv.Atoi(paramId)
	if err != nil || userId <= 0 {
		h.logger.Warnf("Invalid request:%s", err)
		http.Error(w, fmt.Sprintf("Invalid url request:%s", err), 400)
		return
	}
	err = h.service.AppUsers.ChangeUser(&input, userId)
	if err != nil {
		if err.Error() == "such a user does not exist" {
			h.logger.Warnf("changeUser: such user does not exist")
			http.Error(w, "such user does not exist", 404)
			return
		}
		h.logger.Errorf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteUser(w http.ResponseWriter, req *http.Request) {
	paramId := strings.TrimPrefix(req.URL.Path, "/users/")
	userId, err := strconv.Atoi(paramId)
	if err != nil || userId <= 0 {
		h.logger.Warnf("Invalid request:%s", err)
		http.Error(w, fmt.Sprintf("Invalid url request:%s", err), 400)
		return
	}
	err = h.service.AppUsers.DeleteUser(userId)
	if err != nil {
		if err.Error() == "user with such Id does not exist" {
			h.logger.Warnf("deleteUser: such user does not exist")
			http.Error(w, "such user does not exist", 404)
			return
		}
		h.logger.Errorf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getHobbyByUserId(w http.ResponseWriter, req *http.Request) {
	paramId := strings.TrimPrefix(req.URL.Path, "/users/")
	paramId = strings.TrimSuffix(paramId, "/hobbies")
	userId, err := strconv.Atoi(paramId)
	if err != nil || userId <= 0 {
		h.logger.Warnf("Invalid request:%s", err)
		http.Error(w, "Invalid url request", 400)
		return
	}
	hobbiesId, err := h.service.AppUsers.GetHobbyByUserId(userId)
	if err != nil {
		if err.Error() == "user with such Id does not exist" {
			h.logger.Warnf("getHobbyByUserId: such user does not exist")
			http.Error(w, "such user does not exist", 404)
			return
		}
		h.logger.Errorf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(hobbiesId)
	if err != nil {
		h.logger.Errorf("getHobbyByUserId: error while marshaling list id of hobbies: %s", err)
		http.Error(w, fmt.Sprintf("getHobbyByUserId: error while marshaling list id of hobbies: %s", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getHobbyByUserId: error while writing response:%s", err)
		http.Error(w, fmt.Sprintf("getHobbyByUserId: error while writing response:%s", err), 500)
		return
	}
}
