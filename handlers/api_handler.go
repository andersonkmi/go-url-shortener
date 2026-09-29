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

type ApiHandler struct {
	urlShortener *shortener.UrlShortener
}

func NewApiHandler(urlShortener *shortener.UrlShortener) *ApiHandler {
	return &ApiHandler{urlShortener: urlShortener}
}

func (h *ApiHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	if shortCode == "" {
		http.Error(w, "Invalid short code", http.StatusNotFound)
		return
	}

	longUrl, err := h.urlShortener.GetOriginalUrl(shortCode)
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

func (h *ApiHandler) Shorten(writer http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(writer, r.Body, 64*1024)

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

	shortenedUrl, err := h.urlShortener.ShortenUrl(createUrlRequest.URL)
	if err != nil {
		log.WithField("originalURL", createUrlRequest.URL).Warn("Failed to shorten URL")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	shortUrl := fmt.Sprintf("https://%s/%s", r.Host, shortenedUrl)
	generateSuccessResponse(writer, http.StatusCreated, createUrlRequest.URL, shortUrl)
}

func (h *ApiHandler) PerformHealthCheck(writer http.ResponseWriter, _ *http.Request) {
	log.Info("Performing health check")
	healthCheckResponse := HealthCheckResponse{"OK"}
	response, _ := json.Marshal(healthCheckResponse)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write(response); err != nil {
		log.WithError(err).Error("Failed to write health check response")
	}
}

func generateSuccessResponse(writer http.ResponseWriter, code int, originalUrl string, shortenedUrl string) {
	var urlResponse = URLResponse{originalUrl, shortenedUrl}
	response, _ := json.Marshal(urlResponse)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	if _, err := writer.Write(response); err != nil {
		log.WithError(err).Error("Failed to write response")
	}
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
