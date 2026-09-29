package handlers

import (
	"encoding/json"
	"fmt"
	"go-url-shortener/shortener"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// CreateURLRequest is the JSON payload accepted by the Shorten endpoint.
type CreateURLRequest struct {
	URL string `json:"url"`
}

// URLResponse is the JSON body returned after a URL is successfully
// shortened, containing the original URL and its shortened counterpart.
type URLResponse struct {
	URL      string `json:"url"`
	ShortUrl string `json:"shortUrl"`
}

// HealthCheckResponse is the JSON body returned by the health check endpoint.
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// ApiHandler exposes the HTTP handlers for the URL shortener API.
type ApiHandler struct {
	urlShortener *shortener.UrlShortener
}

// NewApiHandler returns an ApiHandler backed by the given UrlShortener.
func NewApiHandler(urlShortener *shortener.UrlShortener) *ApiHandler {
	return &ApiHandler{urlShortener: urlShortener}
}

// Redirect resolves the "shortCode" path value to its original URL and
// replies with a 301 Moved Permanently redirect. It responds with 404 Not
// Found when the short code is missing or unknown, and 500 Internal Server
// Error when the lookup fails.
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

// Shorten reads a CreateURLRequest from the request body (limited to 64 KiB),
// validates the URL, and creates a shortened version of it. On success it
// replies with 201 Created and a URLResponse. It responds with 400 Bad Request
// for a malformed body or invalid URL, and 500 Internal Server Error when
// shortening fails.
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

	shortUrl := fmt.Sprintf("%s://%s/%s", requestScheme(r), r.Host, shortenedUrl)
	generateSuccessResponse(writer, http.StatusCreated, createUrlRequest.URL, shortUrl)
}

// PerformHealthCheck replies with 200 OK and a HealthCheckResponse indicating
// the service is up.
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

// generateSuccessResponse writes a URLResponse as JSON with the given status
// code.
func generateSuccessResponse(writer http.ResponseWriter, code int, originalUrl string, shortenedUrl string) {
	var urlResponse = URLResponse{originalUrl, shortenedUrl}
	response, _ := json.Marshal(urlResponse)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	if _, err := writer.Write(response); err != nil {
		log.WithError(err).Error("Failed to write response")
	}
}

// requestScheme returns the scheme used to reach the service. It honours the
// X-Forwarded-Proto header set by reverse proxies and load balancers, then
// falls back to inspecting whether the request itself was served over TLS.
func requestScheme(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-Proto"); forwarded != "" {
		if scheme, _, found := strings.Cut(forwarded, ","); found {
			return strings.ToLower(strings.TrimSpace(scheme))
		}
		return strings.ToLower(strings.TrimSpace(forwarded))
	}
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

// validateUrl checks that url is non-empty and starts with http:// or
// https://, returning an error describing the first violation found.
func validateUrl(url string) error {
	if url == "" {
		return fmt.Errorf("URL is empty")
	}
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		return fmt.Errorf("URL must start with http or https")
	}
	return nil
}
