package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	repo := repositories.NewRepository(db, logger)
	err = repo.SaveInitialCountries(countries)
	if err != nil {
		logger.Fatal(err)
	}

	ser := services.NewService(repo, logger)
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
	done := make(chan bool, 1)
	go func() {
		if err := serv.Run(host, port, handler.InitRoutes()); err != nil {
			logger.Panicf("Error occured while running http server: %s", err.Error())
			select {
			case <-done:
				return
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				ser.AppCountries.LoadImages()
			}
		}
	}()
	<-quit
	ticker.Stop()
	done <- true
}
