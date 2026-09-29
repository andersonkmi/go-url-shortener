package internal

import (
	"fmt"
	"go-url-shortener/base62"

	log "github.com/sirupsen/logrus"
)

// GenerateShortUrl returns the short code for originalUrl. If the URL was
// shortened before, the existing code is returned; otherwise a new ID is
// drawn from the database sequence, encoded to base62, and persisted.
func (r *ShortUrlRepository) GenerateShortUrl(originalUrl string) (string, error) {
	// Verify if the current URL is already present
	shortenedUrl, shortUrlGetErr := r.getShortenedUrlFromOriginal(originalUrl)
	if shortUrlGetErr != nil {
		log.WithError(shortUrlGetErr).Warn("Failed to retrieve short URL")
		return "", fmt.Errorf("failed to get short url: %w", shortUrlGetErr)
	}

	if shortenedUrl.Url != "" {
		log.WithField("originalURL", originalUrl).Info("URL is already shortened")
		return shortenedUrl.ShortUrl, nil
	}

	id, shortUrlGenErr := r.generateShortUrlId()
	if shortUrlGenErr != nil {
		log.WithError(shortUrlGenErr).Warn("Failed to generate short url id")
		return "", fmt.Errorf("failed to generate short url id: %w", shortUrlGenErr)
	}

	base62Id := base62.IdToBase62(id)
	newShortUrl := ShortUrl{id, originalUrl, base62Id}
	storedShortUrl, urlSaveErr := r.saveShortUrl(newShortUrl)
	if urlSaveErr != nil {
		log.WithField("shortURL", newShortUrl).WithError(urlSaveErr).Warn("Failed to save short URL")
		return "", fmt.Errorf("failed to save short url: %w", urlSaveErr)
	}

	log.WithField("shortURL", newShortUrl).Info("Short URL created")
	return storedShortUrl, nil
}

// GetOriginalUrl looks up the original URL for the given short code. It
// returns an empty string with a nil error when the code is unknown.
func (r *ShortUrlRepository) GetOriginalUrl(shortUrl string) (string, error) {
	url, err := r.getShortenedUrlFromShortenedCode(shortUrl)
	if err != nil {
		log.WithField("shortURL", shortUrl).WithError(err).Warn("Failed to retrieve original URL")
		return "", fmt.Errorf("failed to get original url: %w", err)
	}
	log.WithField("shortURL", shortUrl).Info("Original URL retrieved")
	return url.Url, nil
}
