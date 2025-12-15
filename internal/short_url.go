package internal

import (
	"fmt"
	"go-url-shortener/base62"

	log "github.com/sirupsen/logrus"
)

func GenerateShortUrl(originalUrl string) (string, error) {
	// Verify if the current URL is already present
	shortenedUrl, _ := getShortenedUrlFromOriginal(originalUrl)
	if shortenedUrl.Url != "" {
		log.WithField("originalURL", originalUrl).Info("URL is already shortened")
		return shortenedUrl.ShortUrl, nil
	}

	id, err := generateShortUrlId()
	if err != nil {
		log.WithError(err).Warn("Failed to generate short url id")
		return "", fmt.Errorf("failed to generate short url id: %v", err)
	}

	base62Id := base62.IdToBase62(id)
	newShortUrl := ShortUrl{id, originalUrl, base62Id}
	err2 := saveShortUrl(newShortUrl)
	if err2 != nil {
		log.WithField("shortURL", newShortUrl).WithError(err).Warn("Failed to save short URL")
		return "", fmt.Errorf("failed to save short url: %v", err2)
	}

	log.WithField("shortURL", newShortUrl).Info("Short URL created")
	return base62Id, nil
}

func GetOriginalUrl(shortUrl string) (string, error) {
	url, err := getShortenedUrlFromShortenedCode(shortUrl)
	if err != nil {
		log.WithField("shortURL", shortUrl).WithError(err).Warn("Failed to retrieve original URL")
		return "", fmt.Errorf("failed to get original url: %v", err)
	}
	log.WithField("shortURL", shortUrl).Info("Original URL retrieved")
	return url.Url, nil
}
