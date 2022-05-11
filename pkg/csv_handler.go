package pkg

import (
	"encoding/csv"
	"fmt"
	"lesson_2/pkg/logging"
	"os"
)

var logger = logging.GetLogger()

func CsvHandler(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Errorf("error while opening csv file: %s", err)
		return nil, fmt.Errorf("error while opening csv file: %s", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '	'
	countries, err := reader.ReadAll()
	if err != nil {
		logger.Errorf("error while reading csv file: %s", err)
		return nil, fmt.Errorf("error while reading csv file: %s", err)
	}
	return countries, nil
}
