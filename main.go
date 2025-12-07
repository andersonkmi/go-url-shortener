package main

import (
	"fmt"
	"go-url-shortener/config"
	"go-url-shortener/handlers"
	"go-url-shortener/internal"
	"log"
	"net/http"
)

func main() {
	appConfig := config.LoadConfig()
	if err := internal.Init(appConfig); err != nil {
		log.Fatal(err)
	}
	defer internal.Close()

	fmt.Println("Starting a server on port 8080")
	http.HandleFunc("/", handlers.RedirectHandler)
	http.HandleFunc("/shorten", handlers.ShortenHandler)
	http.HandleFunc("/health", handlers.HealthCheckHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appConfig.ApplicationPort), nil))
}
