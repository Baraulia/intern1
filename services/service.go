package services

import "lesson_2/pkg/logging"

var Countries [][]string

type AppCountries interface {
	GetOneCountry(id string) ([]string, error)
	GetCountries(page int, limit int) ([][]string, int, error)
}

type Service struct {
	AppCountries
}

func NewService(logger logging.Logger) *Service {
	return &Service{
		AppCountries: NewCountryService(logger),
	}
}
