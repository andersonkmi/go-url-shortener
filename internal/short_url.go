package internal

import (
	"fmt"
	"go-url-shortener/base62"
	"log"
)

func GenerateShortUrl(originalUrl string) (string, error) {
	// Verify if the current URL is already present
	shortenedUrl, _ := getShortenedUrlFromOriginal(originalUrl)
	if shortenedUrl.Url != "" {
		log.Default().Println(fmt.Sprintf("URL %s is already shortened", originalUrl))
		return shortenedUrl.ShortUrl, nil
	}

	id, err := generateShortUrlId()
	if err != nil {
		return "", fmt.Errorf("failed to generate short url id: %v", err)
	}

	base62Id := base62.IdToBase62(id)
	newShortUrl := ShortUrl{id, originalUrl, base62Id}
	err2 := saveShortUrl(newShortUrl)
	if err2 != nil {
		return "", fmt.Errorf("failed to save short url: %v", err2)
	}

	return base62Id, nil
}

func GetOriginalUrl(shortUrl string) (string, error) {
	url, err := getShortenedUrlFromShortenedCode(shortUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get original url: %v", err)
	}
	return url.Url, nil
}
