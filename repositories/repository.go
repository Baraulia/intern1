package repositories

import "tranee_service/models"

type AppCountry interface {
	SaveInitialCountries([]models.Country) error
	GetOneCountry(id string) (*models.Country, error)
	GetCountries(filters *models.Filters) ([]models.Country, int, error)
	//GetCountriesWithoutPagination() ([]models.Country, int, error)
	CreateCountry(country *models.ResponseCountry) (string, error)
	ChangeCountry(country *models.ResponseCountry, countryId string) error
	DeleteCountry(countryId string) error
	CheckCountryId(countryId string) error
	//GetCountriesWithoutFlag() ([]models.Country, error)
	LoadImages(countries []models.Country) error
}
