package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting a server on port 8080")
	http.HandleFunc("/", homeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL shortener is up and running")
}
