package services

import "lesson_2/pkg/logging"

type CountryService struct {
	logger logging.Logger
}

func NewCountryService(logger logging.Logger) *CountryService {
	return &CountryService{logger: logger}
}

func (c *CountryService) GetOneCountry(id string) ([]string, error) {
	for _, country := range Countries {
		if country[3] == id || country[4] == id {
			return country, nil
		}
	}
	return []string{}, nil
}

func (c *CountryService) GetCountries(page int, limit int, chunk bool) ([][]string, int, error) {
	if (page == 0 || limit == 0) && chunk == false {
		return Countries, 0, nil
	} else if chunk == false {
		start := 1 + ((page - 1) * limit)
		return Countries[start : start+limit], 0, nil
	}
	return nil, 0, nil
}
