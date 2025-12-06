package internal

import (
	"database/sql"
	"fmt"
)

type ShortUrl struct {
	UrlId    int64
	Url      string
	ShortUrl string
}

func generateShortUrlId(connection *sql.DB) (int64, error) {
	var urlId int64
	err := connection.QueryRow("select nextval('url_id_sequence')").Scan(&urlId)
	if err != nil {
		return -1, fmt.Errorf("failed to generate short url id: %v", err)
	}
	return urlId, nil
}

func saveShortUrl(connection *sql.DB, shortUrl ShortUrl) error {
	_, err := connection.Exec("insert into url(url_id, url, short_url) values ($1, $2, $3)", shortUrl.UrlId, shortUrl.Url, shortUrl.ShortUrl)
	return err
}

func getUrl(connection *sql.DB, url string) (ShortUrl, error) {
	var shortUrl ShortUrl
	err := connection.QueryRow("select url_id, url, short_url from url where url = $1", url).Scan(&shortUrl.UrlId, &shortUrl.Url, &shortUrl.ShortUrl)
	return shortUrl, err
}
