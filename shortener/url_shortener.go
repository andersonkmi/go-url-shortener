// Package shortener implements the URL shortening service layer, delegating
// persistence to the internal repository.
package shortener

import (
	"go-url-shortener/internal"
)

// UrlShortener shortens URLs and resolves short codes back to their
// original URLs.
type UrlShortener struct {
	repo *internal.ShortUrlRepository
}

// New returns a UrlShortener backed by the given repository.
func New(repo *internal.ShortUrlRepository) *UrlShortener {
	return &UrlShortener{repo: repo}
}

// ShortenUrl returns the short code for url, creating and persisting one if
// the URL has not been shortened before.
func (s *UrlShortener) ShortenUrl(url string) (string, error) {
	shortUrl, err := s.repo.GenerateShortUrl(url)
	if err != nil {
		return "", err
	}
	return shortUrl, nil
}

// GetOriginalUrl resolves shortUrl to its original URL. It returns an empty
// string with a nil error when the short code is unknown.
func (s *UrlShortener) GetOriginalUrl(shortUrl string) (string, error) {
	originalUrl, err := s.repo.GetOriginalUrl(shortUrl)
	if err != nil {
		return "", err
	}
	return originalUrl, nil
}
