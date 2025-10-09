package main

import (
	"fmt"
	"go-url-shortener/internal/store"
	"log"
	"net/http"
)

var urlStore *store.URLStore

func main() {
	urlStore = store.NewURLStore()
	fmt.Println("Starting a server on port 8080")
	http.HandleFunc("/", homeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL shortener is up and running")
}
