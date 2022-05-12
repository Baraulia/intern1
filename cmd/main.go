package main

import (
	"github.com/joho/godotenv"
	"lesson_2/handlers"
	"lesson_2/pkg"
	"lesson_2/pkg/logging"
	"lesson_2/pkg/server"
	"lesson_2/services"
	"os"
)

func main() {
	logger := logging.GetLogger()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalf("Error loading .env file. %s", err.Error())
	}

	path := os.Getenv("PATH_CSV_FILE")
	countries, err := pkg.CsvHandler(path)
	if err != nil {
		logger.Fatal(err)
	}
	services.Countries = countries[1:]

	ser := services.NewService(logger)
	handler := handlers.NewHandler(ser, logger)

	port := os.Getenv("API_SERVER_PORT")
	host := os.Getenv("API_SERVER_HOST")
	serv := new(server.Server)

	logger.Infof("Starting server on %s:%s...", host, port)
	if err := serv.Run(host, port, handler.InitRoutes()); err != nil {
		logger.Panicf("Error occured while running http server: %s", err.Error())
	}

}
