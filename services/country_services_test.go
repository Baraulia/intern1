package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"tranee_service/internal/logging"
	"tranee_service/models"
)

func TestService_GetCountries(t *testing.T) {
	testTable := []struct {
		name              string
		inputPage         int
		inputLimit        int
		expectedCountries []models.Country
		expectedError     error
	}{
		{
			name:       "OK",
			inputLimit: 2,
			inputPage:  1,
			expectedCountries: []models.Country{
				{
					Name:            "Абхазия",
					FullName:        "Республика Абхазия",
					EnglishName:     "Abkhazia",
					Alpha2:          "AB",
					Alpha3:          "ABH",
					Iso:             895,
					Location:        "Азия",
					LocationPrecise: "Закавказье",
				},
				{
					Name:            "Австралия",
					FullName:        "",
					EnglishName:     "Australia",
					Alpha2:          "AU",
					Alpha3:          "AUS",
					Iso:             36,
					Location:        "Океания",
					LocationPrecise: "Австралия и Новая Зеландия",
				},
			},
			expectedError: nil,
		},
		{
			name:              "out of range",
			inputLimit:        3,
			inputPage:         2,
			expectedCountries: nil,
			expectedError:     errors.New("limit out of range"),
		},
		{
			name:       "without pagination",
			inputLimit: 0,
			inputPage:  0,
			expectedCountries: []models.Country{
				{
					Name:            "Абхазия",
					FullName:        "Республика Абхазия",
					EnglishName:     "Abkhazia",
					Alpha2:          "AB",
					Alpha3:          "ABH",
					Iso:             895,
					Location:        "Азия",
					LocationPrecise: "Закавказье",
				},
				{
					Name:            "Австралия",
					FullName:        "",
					EnglishName:     "Australia",
					Alpha2:          "AU",
					Alpha3:          "AUS",
					Iso:             36,
					Location:        "Океания",
					LocationPrecise: "Австралия и Новая Зеландия",
				},
				{
					Name:            "Австрия",
					FullName:        "Австрийская Республика",
					EnglishName:     "Austria",
					Alpha2:          "AT",
					Alpha3:          "AUT",
					Iso:             40,
					Location:        "Европа",
					LocationPrecise: "Западная Европа",
				},
				{
					Name:            "Азербайджан",
					FullName:        "Республика Азербайджан",
					EnglishName:     "Azerbaijan",
					Alpha2:          "AZ",
					Alpha3:          "AZE",
					Iso:             31,
					Location:        "Азия",
					LocationPrecise: "Западная Азия",
				},
				{
					Name:            "Албания",
					FullName:        "Республика Албания",
					EnglishName:     "Albania",
					Alpha2:          "AL",
					Alpha3:          "ALB",
					Iso:             8,
					Location:        "Европа",
					LocationPrecise: "Южная Европа",
				},
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			logger := logging.GetLogger()
			countries := [][]string{
				{"Абхазия", "Республика Абхазия", "Abkhazia", "AB", "ABH", "895", "Азия", "Закавказье"},
				{"Австралия", "", "Australia", "AU", "AUS", "036", "Океания", "Австралия и Новая Зеландия"},
				{"Австрия", "Австрийская Республика", "Austria", "AT", "AUT", "040", "Европа", "Западная Европа"},
				{"Азербайджан", "Республика Азербайджан", "Azerbaijan", "AZ", "AZE", "031", "Азия", "Западная Азия"},
				{"Албания", "Республика Албания", "Albania", "AL", "ALB", "008", "Европа", "Южная Европа"},
			}
			service := NewCountryService(countries, logger)
			gotCountries, _, err := service.GetCountries(testCase.inputPage, testCase.inputLimit)

			//Assert
			assert.Equal(t, testCase.expectedCountries, gotCountries)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestService_GetOneCountry(t *testing.T) {
	testTable := []struct {
		name            string
		inputId         string
		expectedCountry *models.Country
		expectedError   error
	}{
		{
			name:    "OK alpha2",
			inputId: "AB",
			expectedCountry: &models.Country{
				Name:            "Абхазия",
				FullName:        "Республика Абхазия",
				EnglishName:     "Abkhazia",
				Alpha2:          "AB",
				Alpha3:          "ABH",
				Iso:             895,
				Location:        "Азия",
				LocationPrecise: "Закавказье",
			},
			expectedError: nil,
		},
		{
			name:    "OK alpha3",
			inputId: "ABH",
			expectedCountry: &models.Country{
				Name:            "Абхазия",
				FullName:        "Республика Абхазия",
				EnglishName:     "Abkhazia",
				Alpha2:          "AB",
				Alpha3:          "ABH",
				Iso:             895,
				Location:        "Азия",
				LocationPrecise: "Закавказье",
			},
			expectedError: nil,
		},
		{
			name:            "does not exist",
			inputId:         "XXX",
			expectedCountry: nil,
			expectedError:   errors.New("such a country does not exist"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			logger := logging.GetLogger()
			countries := [][]string{
				{"Абхазия", "Республика Абхазия", "Abkhazia", "AB", "ABH", "895", "Азия", "Закавказье"},
				{"Австралия", "", "Australia", "AU", "AUS", "036", "Океания", "Австралия и Новая Зеландия"},
				{"Австрия", "Австрийская Республика", "Austria", "AT", "AUT", "040", "Европа", "Западная Европа"},
				{"Азербайджан", "Республика Азербайджан", "Azerbaijan", "AZ", "AZE", "031", "Азия", "Западная Азия"},
				{"Албания", "Республика Албания", "Albania", "AL", "ALB", "008", "Европа", "Южная Европа"},
			}
			service := NewCountryService(countries, logger)
			gotCountry, err := service.GetOneCountry(testCase.inputId)

			//Assert
			assert.Equal(t, testCase.expectedCountry, gotCountry)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
