// Package internal contains the persistence layer of the URL shortener:
// database initialization and the repository for short URL records.
package internal

import (
	"context"
	"database/sql"
	"fmt"
	"go-url-shortener/config"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// InitDB opens a PostgreSQL connection pool using the given configuration,
// applies the pool limits and connection lifetimes from config, and verifies
// connectivity with a ping. The caller is responsible for closing the
// returned handle.
func InitDB(config config.Config) (*sql.DB, error) {
	connectionStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(config.ConnectionMaxLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(config.ConnectionMaxIdleTime) * time.Minute)

	if err := db.PingContext(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connection established")
	return db, nil
}
