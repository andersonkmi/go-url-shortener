package internal

import (
	"database/sql"
	"fmt"
	"go-url-shortener/config"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init(config config.Config) error {
	connectionStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
	var err error
	db, err = sql.Open("postgres", connectionStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(config.ConnectionMaxLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(config.ConnectionMaxIdleTime) * time.Minute)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
