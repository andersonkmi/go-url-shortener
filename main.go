package main

import (
	"encoding/json"
	"fmt"
	"go-url-shortener/internal"
	"go-url-shortener/shortener"
	"go-url-shortener/store"
	"log"
	"net/http"
	"strings"
)

var urlStore *store.URLStore

type CreateURLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	URL      string `json:"url"`
	ShortUrl string `json:"shortUrl"`
}

func main() {
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

	var createUrlRequest CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&createUrlRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if createUrlRequest.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Get a connection to the database
	connection := internal.GetDB()
	shortenedUrl, err := shortener.ShortenUrl(createUrlRequest.URL, connection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortUrl := fmt.Sprintf("http://localhost:8080/%s", shortenedUrl)
	fmt.Sprintln(w, "Short URL created: %s", shortUrl)

	var urlResponse = URLResponse{createUrlRequest.URL, shortUrl}
	respondWithJSON(w, http.StatusOK, urlResponse)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
