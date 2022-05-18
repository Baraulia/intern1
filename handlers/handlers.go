package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"tranee_service/models"
)

func (h *Handler) getAllCountries(w http.ResponseWriter, req *http.Request) {
	var page = 0
	var limit = 0
	var chunk = false

	if req.URL.Query().Get("page") != "" {
		paramPage, err := strconv.Atoi(req.URL.Query().Get("page"))
		if err != nil || paramPage < 0 {
			h.logger.Warnf("Invalid url request:%s", err)
			http.Error(w, fmt.Sprintf("Invalid url request:%s", err), 400)
			return
		}
		page = paramPage
	}
	if req.URL.Query().Get("limit") != "" {
		paramLimit, err := strconv.Atoi(req.URL.Query().Get("limit"))
		if err != nil || paramLimit < 0 {
			h.logger.Warnf("Invalid url request:%s", err)
			http.Error(w, fmt.Sprintf("Invalid url request:%s", err), 400)
			return
		}
		limit = paramLimit
	}
	if req.URL.Query().Get("chunk") != "" {
		paramChunk := req.URL.Query().Get("chunk")
		if paramChunk != "true" && paramChunk != "false" {
			h.logger.Warnf("Invalid parameter 'chunk' passed")
			http.Error(w, fmt.Sprintf("Invalid parameter 'chunk' passed"), 400)
			return
		}
		if paramChunk == "true" {
			chunk = true
		} else {
			chunk = false
		}
	}

	countries, pages, err := h.service.GetCountries(page, limit)
	if err != nil {
		h.logger.Warnf("server error: %s", err)
		http.Error(w, "server error", 500)
		return
	}

	if chunk == false {
		output, err := json.Marshal(countries)
		if err != nil {
			h.logger.Errorf("getAllCountries: error while marshaling list of countries: %s", err)
			http.Error(w, fmt.Sprintf("getAllCountries: error while marshaling list of countries: %s", err), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Pages", strconv.Itoa(pages))
		_, err = w.Write(output)
		if err != nil {
			h.logger.Errorf("getAllCountries: error while writing response:%s", err)
			http.Error(w, fmt.Sprintf("getAllCountries: error while writing response:%s", err), 500)
			return
		}
	} else {
		flusher, ok := w.(http.Flusher)
		if !ok {
			h.logger.Errorf("getAllCountries: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Connection", "keep-alive")
		for _, country := range countries {
			output, err := json.Marshal(country)
			if err != nil {
				h.logger.Errorf("getAllCountries: error while marshaling one country: %s", err)
				http.Error(w, fmt.Sprintf("getAllCountries: error while marshaling marshaling one country: %s", err), 500)
				return
			}
			addBytes := []byte("\n")
			for _, b := range addBytes {
				output = append(output, b)
			}
			_, err = w.Write(output)
			if err != nil {
				h.logger.Errorf("getAllCountries: error while writing response:%s", err)
				http.Error(w, fmt.Sprintf("getAllCountries: error while writing response:%s", err), 500)
				return
			}
			flusher.Flush()
		}
	}
}

func (h *Handler) getOneCountry(w http.ResponseWriter, req *http.Request) {
	countryId := strings.TrimPrefix(req.URL.Path, "/countries/")
	countryId = strings.ToUpper(countryId)

	country, err := h.service.GetOneCountry(countryId)
	if err != nil {
		if err.Error() == "such a country does not exist" {
			h.logger.Warnf("getOneCountry: such country does not exist")
			http.Error(w, "such country does not exist", 404)
			return
		}
		h.logger.Warnf("getOneCountry: server error: %s", err)
		http.Error(w, "server error", 500)
		return
	}
	output, err := json.Marshal(country)
	if err != nil {
		h.logger.Errorf("getOneCountry: error while marshaling one country: %s", err)
		http.Error(w, fmt.Sprintf("getOneCountry: error while marshaling one country: %s", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Errorf("getOneCountry: error while writing response:%s", err)
		http.Error(w, fmt.Sprintf("getOneCountry: error while writing response:%s", err), 500)
		return
	}
}

func (h *Handler) createCountry(w http.ResponseWriter, req *http.Request) {
	var input models.ResponseCountry
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Errorf("Error while decoding request:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	validationErrors := validateStruct(input)
	if len(validationErrors) != 0 {
		h.logger.Warnf("Incorrect data came from the request:%s", validationErrors)
		errors, err := json.Marshal(validationErrors)
		if err != nil {
			h.logger.Errorf("CreateCountry: error while marshaling list myErrors:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errors)
		if err != nil {
			h.logger.Errorf("CreateCountry: can not write errors into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	countryId, err := h.service.CreateCountry(&input)
	if err != nil {
		if err != nil {
			h.logger.Errorf(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
	}
	w.Header().Set("id", countryId)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) changeCountry(w http.ResponseWriter, req *http.Request) {
	var input models.ResponseCountry
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Errorf("Error while decoding request:%s", err)
		http.Error(w, err.Error(), 400)
		return
	}
	validationErrors := validateStruct(input)
	if len(validationErrors) != 0 {
		h.logger.Warnf("Incorrect data came from the request:%s", validationErrors)
		errors, err := json.Marshal(validationErrors)
		if err != nil {
			h.logger.Errorf("ChangeCountry: error while marshaling list myErrors:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errors)
		if err != nil {
			h.logger.Errorf("ChangeCountry: can not write errors into response:%s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	countryId := strings.TrimPrefix(req.URL.Path, "/countries/")
	countryId = strings.ToUpper(countryId)
	err := h.service.ChangeCountry(&input, countryId)
	if err != nil {
		if err.Error() == "such a country does not exist" {
			h.logger.Warnf("changeCountry: such country does not exist")
			http.Error(w, "such country does not exist", 404)
			return
		}
		h.logger.Errorf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteCountry(w http.ResponseWriter, req *http.Request) {
	reqId := strings.TrimPrefix(req.URL.Path, "/countries/")
	reqId = strings.ToUpper(reqId)
	err := h.service.DeleteCountry(reqId)
	if err != nil {
		if err.Error() == "such a country does not exist" {
			h.logger.Warnf("deleteCountry: such country does not exist")
			http.Error(w, "such country does not exist", 404)
			return
		}
		h.logger.Errorf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) loadImages(w http.ResponseWriter, req *http.Request) {
	go h.service.LoadImages()
	go w.WriteHeader(http.StatusProcessing)
}
