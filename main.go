package main

import (
	"context"
	"errors"
	"fmt"
	"go-url-shortener/config"
	"go-url-shortener/handlers"
	"go-url-shortener/internal"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	appConfig := config.LoadConfig()
	if err := internal.InitDB(appConfig); err != nil {
		log.Fatal(err)
	}

	log.Info("Starting application...")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{shortCode}", handlers.RedirectHandler)
	mux.HandleFunc("POST /shorten", handlers.ShortenHandler)
	mux.HandleFunc("GET /health", handlers.HealthCheckHandler)
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", appConfig.ApplicationPort),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Fatal("Server error")
		}
	}()

	log.Info("Server started")

	// Block until signal received
	<-sigChan
	log.Info("Shutting down...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server gracefully
	if err := httpServer.Shutdown(ctx); err != nil {
		log.WithError(err).Warn("Server shutdown error")
	}

	// Close database connection
	if err := internal.CloseDB(); err != nil {
		log.WithError(err).Warn("Database close error")
	}
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000Z07:00",
	})
}
