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
		return -1, fmt.Errorf("failed to generate short url id: %v", err)
	}
	return urlId, nil
}

func saveShortUrl(shortUrl ShortUrl) error {
	_, err := db.Exec("insert into url(url_id, url, short_url) values ($1, $2, $3)", shortUrl.UrlId, shortUrl.Url, shortUrl.ShortUrl)
	return err
}

func getShortenedUrlFromOriginal(url string) (ShortUrl, error) {
	result := db.QueryRow("select url_id, url, short_url from url where url = $1", url)
	if errors.Is(result.Err(), sql.ErrNoRows) {
		emptyResult := ShortUrl{0, "", ""}
		return emptyResult, nil
	}

	// Returns a valid result
	var shortUrl ShortUrl
	err := result.Scan(&shortUrl.UrlId, &shortUrl.Url, &shortUrl.ShortUrl)
	return shortUrl, err
}

func getShortenedUrlFromShortenedCode(shortenedCode string) (ShortUrl, error) {
	result := db.QueryRow("select url_id, url, short_url from url where short_url = $1", shortenedCode)
	if errors.Is(result.Err(), sql.ErrNoRows) {
		emptyResult := ShortUrl{0, "", ""}
		return emptyResult, nil
	}

	// Returns a valid result
	var shortUrl ShortUrl
	err := result.Scan(&shortUrl.UrlId, &shortUrl.Url, &shortUrl.ShortUrl)
	return shortUrl, err
}
