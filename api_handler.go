package main

import (
	"encoding/json"
	"fmt"
	"go-url-shortener/shortener"
	"log"
	"net/http"
	"strings"
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Remove landing slash from path
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" {
		// should return something to the API caller
		return
	}

	// Replace with the new function
	longUrl := ""
	http.Redirect(w, r, longUrl, http.StatusFound)
}

func shortenHandler(writer http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var createUrlRequest CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&createUrlRequest); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := validateUrl(createUrlRequest.URL)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	shortenedUrl, err := shortener.ShortenUrl(createUrlRequest.URL)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	shortUrl := fmt.Sprintf("http://localhost:8080/%s", shortenedUrl)
	fmt.Sprintln(writer, "Short URL created: %s", shortUrl)

	generateSuccessResponse(writer, http.StatusCreated, createUrlRequest.URL, shortUrl)
}

func generateSuccessResponse(writer http.ResponseWriter, code int, originalUrl string, shortenedUrl string) {
	var urlResponse = URLResponse{originalUrl, shortenedUrl}
	response, _ := json.Marshal(urlResponse)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(response)
}

func validateUrl(url string) error {
	if url == "" {
		return fmt.Errorf("URL is empty")
	}
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		return fmt.Errorf("URL must start with http or https")
	}
	return nil
}
