package shortener

import (
	"go-url-shortener/internal"
)

func ShortenUrl(url string) (string, error) {
	connection := internal.GetDB()
	shortUrl, err := internal.GenerateShortUrl(url, connection)
	if err != nil {
		return "", err
	}
	return shortUrl, nil
}
