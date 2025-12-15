package handlers

import (
	"encoding/json"
	"fmt"
	"go-url-shortener/shortener"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type CreateURLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	URL      string `json:"url"`
	ShortUrl string `json:"shortUrl"`
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" {
		http.Error(w, "Invalid short code", http.StatusNotFound)
		return
	}

	longUrl, err := shortener.GetOriginalUrl(shortCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if longUrl == "" {
		http.Error(w, "Invalid short code", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longUrl, http.StatusMovedPermanently)
}

func ShortenHandler(writer http.ResponseWriter, r *http.Request) {
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
		log.WithField("originalURL", createUrlRequest.URL).Warn("Invalid URL provided")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	shortenedUrl, err := shortener.ShortenUrl(createUrlRequest.URL)
	if err != nil {
		log.WithField("originalURL", createUrlRequest.URL).Warn("Failed to shorten URL")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	shortUrl := fmt.Sprintf("https://%s/%s", r.Host, shortenedUrl)
	fmt.Sprintln(writer, "Short URL created: %s", shortUrl)

	generateSuccessResponse(writer, http.StatusCreated, createUrlRequest.URL, shortUrl)
}

func HealthCheckHandler(writer http.ResponseWriter, r *http.Request) {
	log.Info("Performing health check")
	healthCheckResponse := HealthCheckResponse{"OK"}
	response, _ := json.Marshal(healthCheckResponse)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
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
