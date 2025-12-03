package internal

import (
	"database/sql"
	"fmt"
	"go-url-shortener/base62"
)

func GenerateShortUrl(originalUrl string, connection *sql.DB) (string, error) {
	id, err := generateShortUrlId(connection)
	if err != nil {
		return "", fmt.Errorf("failed to generate short url id: %v", err)
	}

	base62Id := base62.IdToBase62(id)
	shortUrl := ShortUrl{id, originalUrl, base62Id}
	err2 := saveShortUrl(connection, shortUrl)
	if err2 != nil {
		return "", fmt.Errorf("failed to save short url: %v", err2)
	}

	return base62Id, nil
}
