package services

import "tranee_service/models"

type AppCountries interface {
	GetOneCountry(id string) (*models.Country, error)
	GetCountries(page int, limit int) ([]models.Country, int, error)
	CreateCountry(country *models.ResponseCountry) (string, error)
	ChangeCountry(country *models.ResponseCountry, countryId string) error
	DeleteCountry(countryId string) error
	LoadImages()
}
