package shortener

import (
	"go-url-shortener/internal"
)

func ShortenUrl(url string) (string, error) {
	shortUrl, err := internal.GenerateShortUrl(url)
	if err != nil {
		return "", err
	}
	return shortUrl, nil
}
