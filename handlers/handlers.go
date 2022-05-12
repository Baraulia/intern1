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

	countries, pages, err := h.service.AppCountries.GetCountries(page, limit)
	if err != nil {
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
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")
		for _, country := range countries {
			output, err := json.Marshal(country)
			if err != nil {
				h.logger.Errorf("getAllCountries: error while marshaling one country: %s", err)
				http.Error(w, fmt.Sprintf("getAllCountries: error while marshaling marshaling one country: %s", err), 500)
				return
			}
			addBytes := []byte("\r\n")
			for _, b := range addBytes {
				output = append(output, b)
			}
			w.Header().Set("Content-Length", strconv.Itoa(len(output))+"\r\n")
			w.Header().Set("Transfer-Encoding", "chunked")
			_, err = w.Write([]byte(strconv.Itoa(len(output)) + "\r\n"))
			if err != nil {
				h.logger.Errorf("getAllCountries: error while writing response:%s", err)
				http.Error(w, fmt.Sprintf("getAllCountries: error while writing response:%s", err), 500)
				return
			}
			_, err = w.Write(output)
			if err != nil {
				h.logger.Errorf("getAllCountries: error while writing response:%s", err)
				http.Error(w, fmt.Sprintf("getAllCountries: error while writing response:%s", err), 500)
				return
			}
		}
	}
}

func (h *Handler) getOneCountry(w http.ResponseWriter, req *http.Request) {
	countryId := strings.TrimPrefix(req.URL.Path, "/countries/")

	country, err := h.service.AppCountries.GetOneCountry(countryId)
	if err != nil {
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
