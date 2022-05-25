package repositories

import (
	"database/sql"
	"tranee_service/internal/logging"
	"tranee_service/models"
)

type AppCountry interface {
	SaveInitialCountries([]models.Country) error
	GetOneCountry(id string) (*models.Country, error)
	GetCountries(filters *models.Filters) ([]models.Country, int, error)
	CreateCountry(country *models.ResponseCountry) (string, error)
	ChangeCountry(country *models.ResponseCountry, countryId string) error
	DeleteCountry(countryId string) error
	CheckCountryId(countryId string) error
	LoadImages(countries []models.Country) error
}

type AppUsers interface {
	CreateUser(user *models.User) (int, error)
	GetUserById(userId int) (*models.ResponseUser, error)
	GetUsers(filter *models.Options) ([]models.ResponseUser, int, error)
	ChangeUser(user *models.User, userId int) error
	DeleteUser(userId int) error
}

type Repository struct {
	AppCountry
	AppUsers
}

func NewRepository(db *sql.DB, logger logging.Logger) *Repository {
	return &Repository{
		AppCountry: NewCountryRepository(db, logger),
		AppUsers:   NewUserRepository(db, logger),
	}
}
