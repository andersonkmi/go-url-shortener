package shortener

import (
	"go-url-shortener/internal"
)

type UrlShortener struct {
	repo *internal.ShortUrlRepository
}

func New(repo *internal.ShortUrlRepository) *UrlShortener {
	return &UrlShortener{repo: repo}
}

func (s *UrlShortener) ShortenUrl(url string) (string, error) {
	shortUrl, err := s.repo.GenerateShortUrl(url)
	if err != nil {
		return "", err
	}
	return shortUrl, nil
}

func (s *UrlShortener) GetOriginalUrl(shortUrl string) (string, error) {
	originalUrl, err := s.repo.GetOriginalUrl(shortUrl)
	if err != nil {
		return "", err
	}
	return originalUrl, nil
}
