package services

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"tranee_service/internal/logging"
	"tranee_service/models"
	"tranee_service/repositories"
)

type CountryService struct {
	repository repositories.AppCountry
	logger     logging.Logger
}

func NewCountryService(repository repositories.AppCountry, logger logging.Logger) *CountryService {
	return &CountryService{repository: repository, logger: logger}
}

func (c *CountryService) GetOneCountry(id string) (*models.Country, error) {
	if err := c.repository.CheckCountryId(id); err != nil {
		return nil, err
	}
	return c.repository.GetOneCountry(id)
}

func (c *CountryService) GetCountries(filters *models.Filters) ([]models.Country, int, error) {
	return c.repository.GetCountries(filters)
}

func (c *CountryService) CreateCountry(country *models.ResponseCountry) (string, error) {
	return c.repository.CreateCountry(country)
}

func (c *CountryService) ChangeCountry(country *models.ResponseCountry, countryId string) error {
	if err := c.repository.CheckCountryId(countryId); err != nil {
		return err
	}
	return c.repository.ChangeCountry(country, countryId)
}

func (c *CountryService) DeleteCountry(countryId string) error {
	return c.repository.DeleteCountry(countryId)
}

func (c *CountryService) LoadImages() {
	countries, _, err := c.repository.GetCountries(&models.Filters{
		Page:  0,
		Limit: 0,
		Flag:  true,
	})
	if err != nil {
		c.logger.Errorf(err.Error())
	}
	var changedCountries []models.Country
	for _, country := range countries {
		request := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&prop=pageimages&format=json&formatversion=2&piprop=original&titles=%s", country.EnglishName)
		response, err := http.Get(request)
		if err != nil {
			c.logger.Errorf("Error while sending request to wikipedia:%s", err)
		}
		defer response.Body.Close()
		b, err := io.ReadAll(response.Body)
		if err != nil {
			c.logger.Errorf("Error while sending request to wikipedia:%s", err)
		}
		url := gjson.Get(string(b), "query.pages.0.original.source")
		if url.String() != "" {
			country.Url = url.String()
			changedCountries = append(changedCountries, country)
		}
	}
	err = c.repository.LoadImages(changedCountries)
	if err != nil {
		c.logger.Errorf("Error while saving images url:%s", err)
	}
}
