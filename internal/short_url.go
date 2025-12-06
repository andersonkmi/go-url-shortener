package internal

import (
	"database/sql"
	"fmt"
	"go-url-shortener/base62"
	"log"
)

func GenerateShortUrl(originalUrl string, connection *sql.DB) (string, error) {
	// Verify if the current URL is already present
	shortUrl, err := getUrl(connection, originalUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get short url: %v", err)
	}

	if shortUrl.Url != "" {
		log.Default().Println(fmt.Sprintf("URL %s is already shortened", originalUrl))
		return shortUrl.ShortUrl, nil
	}

	id, err := generateShortUrlId(connection)
	if err != nil {
		return "", fmt.Errorf("failed to generate short url id: %v", err)
	}

	base62Id := base62.IdToBase62(id)
	newShortUrl := ShortUrl{id, originalUrl, base62Id}
	err2 := saveShortUrl(connection, newShortUrl)
	if err2 != nil {
		return "", fmt.Errorf("failed to save short url: %v", err2)
	}

	return base62Id, nil
}
