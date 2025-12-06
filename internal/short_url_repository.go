package internal

import (
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

func getUrl(url string) (ShortUrl, error) {
	var shortUrl ShortUrl
	err := db.QueryRow("select url_id, url, short_url from url where url = $1", url).Scan(&shortUrl.UrlId, &shortUrl.Url, &shortUrl.ShortUrl)
	return shortUrl, err
}

func urlAlreadyExists(url string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM url WHERE url = $1)`
	db.QueryRow(query, url).Scan(&exists)
	return exists
}
