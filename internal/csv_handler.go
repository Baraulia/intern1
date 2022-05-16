package internal

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func CsvHandler(filePath, separator string) ([][]string, error) {
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
