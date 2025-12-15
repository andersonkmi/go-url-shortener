package internal

import (
	"database/sql"
	"fmt"
	"go-url-shortener/config"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

func InitDB(config config.Config) error {
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
		log.WithError(err).Fatal("Failed to ping database")
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connection established")
	return nil
}

func CloseDB() error {
	if db != nil {
		log.Info("Closing database connection")
		return db.Close()
	}
	return nil
}
