package repository

import (
	"database/sql"
	"fmt"
)

type ShortUrl struct {
	UrlId    int64
	Url      string
	ShortUrl string
}

func GenerateShortUrlId(connection *sql.DB) (int64, error) {
	var urlId int64
	err := connection.QueryRow("select nextval('url_id_sequence')").Scan(&urlId)
	if err != nil {
		return -1, fmt.Errorf("Failed to generate short url id: %v", err)
	}
	return urlId, nil
}
