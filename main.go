package main

import (
	"fmt"
	"go-url-shortener/config"
	"go-url-shortener/internal"
	"log"
	"net/http"
)

func main() {
	dbConfig := config.LoadDBConfig()
	if err := internal.Init(dbConfig); err != nil {
		log.Fatal(err)
	}
	defer internal.Close()

	fmt.Println("Starting a server on port 8080")
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/shorten", shortenHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
