package internal

import (
	"fmt"
	"go-url-shortener/base62"

	log "github.com/sirupsen/logrus"
)

func GenerateShortUrl(originalUrl string) (string, error) {
	// Verify if the current URL is already present
	shortenedUrl, shortUrlGetErr := getShortenedUrlFromOriginal(originalUrl)
	if shortUrlGetErr != nil {
		log.WithError(shortUrlGetErr).Warn("Failed to retrieve short URL")
		return "", fmt.Errorf("failed to get short url: %v", shortUrlGetErr)
	}

	if shortenedUrl.Url != "" {
		log.WithField("originalURL", originalUrl).Info("URL is already shortened")
		return shortenedUrl.ShortUrl, nil
	}

	id, shortUrlGenErr := generateShortUrlId()
	if shortUrlGenErr != nil {
		log.WithError(shortUrlGenErr).Warn("Failed to generate short url id")
		return "", fmt.Errorf("failed to generate short url id: %v", shortUrlGenErr)
	}

	base62Id := base62.IdToBase62(id)
	newShortUrl := ShortUrl{id, originalUrl, base62Id}
	storedShortUrl, urlSaveErr := saveShortUrl(newShortUrl)
	if urlSaveErr != nil {
		log.WithField("shortURL", newShortUrl).WithError(urlSaveErr).Warn("Failed to save short URL")
		return "", fmt.Errorf("failed to save short url: %v", urlSaveErr)
	}

	log.WithField("shortURL", newShortUrl).Info("Short URL created")
	return storedShortUrl, nil
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
