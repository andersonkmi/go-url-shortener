package main

import (
	"fmt"
	"go-url-shortener/config"
	"go-url-shortener/handlers"
	"go-url-shortener/internal"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	appConfig := config.LoadConfig()
	if err := internal.InitDB(appConfig); err != nil {
		log.Fatal(err)
	}
	defer internal.CloseDB()

	log.Info("Starting application...")
	http.HandleFunc("/", handlers.RedirectHandler)
	http.HandleFunc("/shorten", handlers.ShortenHandler)
	http.HandleFunc("/health", handlers.HealthCheckHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appConfig.ApplicationPort), nil))
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000Z07:00",
	})
}
