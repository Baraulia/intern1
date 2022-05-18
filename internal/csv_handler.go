package internal

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tranee_service/models"
)

func GetCountriesInString(filePath, separator string) ([][]string, error) {
	var responseCountries [][]string
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while opening csv file: %s", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = rune(separator[0])
	countries, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error while reading csv file: %s", err)
	}
	countries = countries[1:]
	for _, country := range countries {
		responseCountry := strings.Split(country[0], "	")
		responseCountries = append(responseCountries, responseCountry)
	}
	return responseCountries, nil
}

func CsvHandler(filePath, separator string) ([]models.Country, error) {
	var countries []models.Country
	stringCountries, err := GetCountriesInString(filePath, separator)
	if err != nil {
		return nil, err
	}
	for _, country := range stringCountries {
		var countryStruct models.Country
		countryStruct.Name = country[0]
		countryStruct.FullName = country[1]
		countryStruct.EnglishName = country[2]
		countryStruct.Alpha2 = country[3]
		countryStruct.Alpha3 = country[4]
		countryStruct.Iso, _ = strconv.Atoi(country[5])
		countryStruct.Location = country[6]
		countryStruct.LocationPrecise = country[7]
		countries = append(countries, countryStruct)
	}
	return countries, nil
}
