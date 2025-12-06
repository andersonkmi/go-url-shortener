package main

import (
	"fmt"
	"go-url-shortener/internal"
	"log"
	"net/http"
)

func main() {
	// Change this to use parameters from a config function
	connectionStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		"127.0.0.1", 5432, "pguser", "pgpwd", "urlshortner", "disable")
	if err := internal.Init(connectionStr); err != nil {
		log.Fatal(err)
	}
	defer internal.Close()

	fmt.Println("Starting a server on port 8080")
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/shorten", shortenHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
