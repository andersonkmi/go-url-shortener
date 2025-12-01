package main

import (
	"fmt"
	"go-url-shortener/shortener"
	"go-url-shortener/store"
	"log"
	"net/http"
	"strings"
)

var urlStore *store.URLStore

func main() {
	urlStore = store.NewURLStore()
	fmt.Println("Starting a server on port 8080")
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/shorten", shortenHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Remove landing slash from path
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" {
		fmt.Fprintf(w, "URL shortener is up and running.\n\nUsage:\n")
		fmt.Fprintf(w, "POST /shorten with url=<originalUrl> to create a short link\n")
		return
	}

	// Look up the URL
	longUrl, exists := urlStore.Get(shortCode)
	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longUrl, http.StatusFound)
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
	fmt.Sprintln(w, "Short URL created: %s", shortUrl)
}
