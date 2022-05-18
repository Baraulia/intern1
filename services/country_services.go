package services

import (
	"fmt"
	"tranee_service/models"
	"tranee_service/repositories"
)

type CountryService struct {
	repository repositories.AppCountry
	logger     models.Logger
}

func NewCountryService(repository repositories.AppCountry, logger models.Logger) *CountryService {
	return &CountryService{repository: repository, logger: logger}
}

func (c *CountryService) GetOneCountry(id string) (*models.Country, error) {
	if err := c.repository.CheckCountryId(id); err != nil {
		return nil, err
	}
	return c.repository.GetOneCountry(id)
}

func (c *CountryService) GetCountries(page int, limit int) ([]models.Country, int, error) {
	if page == 0 || limit == 0 {
		return c.repository.GetCountriesWithoutPagination()
	} else {
		return c.repository.GetCountries(page, limit)
	}
}

func (c *CountryService) CreateCountry(country *models.ResponseCountry) (string, error) {
	return c.repository.CreateCountry(country)
}

func (c *CountryService) ChangeCountry(country *models.ResponseCountry, countryId string) error {
	if err := c.repository.CheckCountryId(countryId); err != nil {
		return err
	}
	if len(countryId) == 2 {
		country.Alpha2 = countryId
	} else if len(countryId) == 3 {
		country.Alpha3 = countryId
	} else {
		return fmt.Errorf("incorrect id")
	}

	return c.repository.ChangeCountry(country, countryId)
}

func (c *CountryService) DeleteCountry(countryId string) error {
	if err := c.repository.CheckCountryId(countryId); err != nil {
		return err
	}
	return c.repository.DeleteCountry(countryId)
}

//
//func (c *CountryService) LoadImages() error {
//	countries, _, err := c.repository.GetCountriesWithoutPagination()
//	if err != nil {
//		return err
//	}
//	var changedCountries []models.Country
//	for _, country := range countries {
//		go func(country *models.Country) {
//
//		}()
//	}
//	return c.repository.LoadImages()
//}
