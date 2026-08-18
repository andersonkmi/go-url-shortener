package internal

import (
	"database/sql"
	"errors"
	"fmt"
)

type ShortUrl struct {
	UrlId    int64
	Url      string
	ShortUrl string
}

func generateShortUrlId() (int64, error) {
	var urlId int64
	err := db.QueryRow("select nextval('url_id_sequence')").Scan(&urlId)
	if err != nil {
		return -1, fmt.Errorf("failed to generate short url id: %w", err)
	}
	return urlId, nil
}

func saveShortUrl(shortUrl ShortUrl) (string, error) {
	var storedShortUrl string
	err := db.QueryRow(
		"insert into url(url_id, url, short_url) values ($1, $2, $3) ON CONFLICT (url) DO UPDATE SET url = EXCLUDED.url RETURNING short_url",
		shortUrl.UrlId, shortUrl.Url, shortUrl.ShortUrl).Scan(&storedShortUrl)
	if err != nil {
		return "", err
	}
	return storedShortUrl, nil
}

func getShortenedUrlFromOriginal(url string) (ShortUrl, error) {
	result := db.QueryRow("select url_id, url, short_url from url where url = $1", url)

	var shortUrl ShortUrl
	err := result.Scan(&shortUrl.UrlId, &shortUrl.Url, &shortUrl.ShortUrl)
	if errors.Is(err, sql.ErrNoRows) {
		emptyResult := ShortUrl{0, "", ""}
		return emptyResult, nil
	}

	// Returns a valid result
	return shortUrl, err
}

func getShortenedUrlFromShortenedCode(shortenedCode string) (ShortUrl, error) {
	result := db.QueryRow("select url_id, url, short_url from url where short_url = $1", shortenedCode)

	var shortUrl ShortUrl
	err := result.Scan(&shortUrl.UrlId, &shortUrl.Url, &shortUrl.ShortUrl)
	if errors.Is(err, sql.ErrNoRows) {
		emptyResult := ShortUrl{0, "", ""}
		return emptyResult, nil
	}

	// Returns a valid result
	return shortUrl, err
}
