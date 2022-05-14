package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
		if err.Error() == "limit out of range" {
			h.logger.Warnf("limit out of range")
			http.Error(w, "limit out of range", 404)
			return
		} else {
			h.logger.Warnf("server error: %s", err)
			http.Error(w, "server error", 500)
			return
		}
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
