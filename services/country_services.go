package services

import (
	"errors"
	"strconv"
	"tranee_service/models"
)

type CountryService struct {
	countries [][]string
	logger    Logger
}

func NewCountryService(counties [][]string, logger Logger) *CountryService {
	return &CountryService{countries: counties, logger: logger}
}

func (c *CountryService) GetOneCountry(id string) (*models.Country, error) {
	var responseCountry models.Country
	exist := false
	for _, country := range c.countries {
		if country[3] == id || country[4] == id {
			responseCountry.Name = country[0]
			responseCountry.FullName = country[1]
			responseCountry.EnglishName = country[2]
			responseCountry.Alpha2 = country[3]
			responseCountry.Alpha3 = country[4]
			responseCountry.Iso, _ = strconv.Atoi(country[5])
			responseCountry.Location = country[6]
			responseCountry.LocationPrecise = country[7]
			exist = true
			return &responseCountry, nil
		}
	}
	if exist == false {
		return nil, errors.New("such a country does not exist")
	}
	return nil, nil
}

func (c *CountryService) GetCountries(page int, limit int) ([]models.Country, int, error) {
	var countries []models.Country
	if page == 0 || limit == 0 {
		for _, country := range c.countries {
			var countryStruct models.Country
			countryStruct.Name = country[0]
			countryStruct.FullName = country[1]
			countryStruct.EnglishName = country[2]
			countryStruct.Alpha2 = country[3]
			countryStruct.Alpha3 = country[4]
			countryStruct.Iso, _ = strconv.Atoi(country[5])
			countryStruct.Location = country[6]
			countryStruct.LocationPrecise = country[7]
			countries = append(countries, countryStruct)
		}
		return countries, 1, nil
	} else {
		start := (page - 1) * limit
		pages := (len(c.countries) - 1) / limit
		if (len(c.countries)-1)%limit != 0 {
			pages++
		}
		if page*limit > len(c.countries) {
			return nil, 0, errors.New("limit out of range")
		}

		for _, country := range c.countries[start : start+limit] {
			var countryStruct models.Country
			countryStruct.Name = country[0]
			countryStruct.FullName = country[1]
			countryStruct.EnglishName = country[2]
			countryStruct.Alpha2 = country[3]
			countryStruct.Alpha3 = country[4]
			countryStruct.Iso, _ = strconv.Atoi(country[5])
			countryStruct.Location = country[6]
			countryStruct.LocationPrecise = country[7]
			countries = append(countries, countryStruct)
		}
		return countries, pages, nil
	}
}
