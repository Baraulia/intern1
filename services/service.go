package services

import "tranee_service/models"

type AppCountries interface {
	GetOneCountry(id string) (*models.Country, error)
	GetCountries(page int, limit int) ([]models.Country, int, error)
}

type Logger interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}
