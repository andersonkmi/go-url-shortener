package shortener

import (
	"database/sql"
	"go-url-shortener/internal"
)

func ShortenUrl(url string, connection *sql.DB) (string, error) {
	shortUrl, err := internal.GenerateShortUrl(url, connection)
	if err != nil {
		return "", err
	}
	return shortUrl, nil
}
