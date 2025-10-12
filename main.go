package main

import (
	"fmt"
	"go-url-shortener/internal/shortener"
	"go-url-shortener/internal/store"
	"log"
	"net/http"
)

var urlStore *store.URLStore

func main() {
	urlStore = store.NewURLStore()
	fmt.Println("Starting a server on port 8080")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/shorten", shortenHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL shortener is up and running")
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	longUrl := r.FormValue("url")
	if longUrl == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortCode := shortener.GenerateShortCode()
	urlStore.Save(shortCode, longUrl)

	shortUrl := fmt.Sprintf("http://localhost:8080/%s", shortCode)
	fmt.Fprintln(w, "Short URL created: %s", shortUrl)
}
