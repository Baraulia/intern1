package services

import (
	"tranee_service/internal/logging"
	"tranee_service/models"
	"tranee_service/repositories"
)

type AppCountries interface {
	GetOneCountry(id string) (*models.Country, error)
	GetCountries(filters *models.Filters) ([]models.Country, int, error)
	CreateCountry(country *models.ResponseCountry) (string, error)
	ChangeCountry(country *models.ResponseCountry, countryId string) error
	DeleteCountry(countryId string) error
	LoadImages()
}

type AppUsers interface {
	CreateUser(user *models.User) (int, error)
	GetUserById(userId int) (*models.ResponseUser, error)
	GetUsers(filter *models.Options) ([]models.ResponseUser, int, error)
	ChangeUser(user *models.User, userId int) error
	DeleteUser(userId int) error
}

type Service struct {
	AppCountries
	AppUsers
}

func NewService(repository *repositories.Repository, logger logging.Logger) *Service {
	return &Service{
		AppCountries: NewCountryService(repository, logger),
		AppUsers:     NewUserService(repository, logger),
	}
}
