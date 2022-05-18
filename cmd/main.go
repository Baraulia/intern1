package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"tranee_service/handlers"
	"tranee_service/internal"
	"tranee_service/internal/databases"
	"tranee_service/internal/logging"
	"tranee_service/internal/server"
	"tranee_service/repositories"
	"tranee_service/services"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file. %s", err.Error())
	}

	separator, present := os.LookupEnv("CSV_SEPARATOR")
	if !present {
		separator = "\t"
	}

	path := os.Getenv("PATH_CSV_FILE")
	countries, err := internal.CsvHandler(path, separator)
	if err != nil {
		log.Fatal(err)
	}

	db, err := databases.NewMysqlDB(&databases.MysqlDB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		Username: os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DBName:   os.Getenv("MYSQL_DATABASE"),
	})
	if err != nil {
		log.Panicf("Error while initialization database:%s", err)
	}
	logger := logging.GetLoggerZap(db)
	repo := repositories.NewCountryRepository(db, logger)
	err = repo.SaveInitialCountries(countries)
	if err != nil {
		logger.Fatal(err)
	}

	ser := services.NewCountryService(repo, logger)
	handler := handlers.NewHandler(ser, logger)

	port, present := os.LookupEnv("API_SERVER_PORT")
	if !present || port == "" {
		port = "0.0.0.0"
	}

	host, present := os.LookupEnv("API_SERVER_HOST")
	if !present || host == "" {
		host = "8090"
	}

	serv := new(server.Server)

	logger.Infof("Starting server on %s:%s...", host, port)
	if err := serv.Run(host, port, handler.InitRoutes()); err != nil {
		logger.Panicf("Error occured while running http server: %s", err.Error())
	}

}
