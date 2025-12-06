package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	DBName                string
	SSLMode               string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime int
	ConnectionMaxIdleTime int
	ApplicationPort       int
}

const (
	dbHostEnvKey                = "DB_HOST"
	dbPortEnvKey                = "DB_PORT"
	dbUserEnvKey                = "DB_USER"
	dbPasswordEnvKey            = "DB_PASSWORD"
	dbNameEnvKey                = "DB_NAME"
	dbSSLModeEnvKey             = "DB_SSL_MODE"
	maxOpenConnectionsEnvKey    = "DB_MAX_OPEN_CONNECTIONS"
	maxIdleConnectionsEnvKey    = "DB_MAX_IDLE_CONNECTIONS"
	connectionMaxLifetimeEnvKey = "DB_CONN_MAX_LIFETIME_MIN"
	connectionMaxIdleEnvKey     = "DB_CONN_MAX_IDLE_TIME_MIN"
	applicationPortEnvKey       = "PORT"
)

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	dbHost := getEnvironment(dbHostEnvKey, "localhost")
	dbPort := getEnvironmentAsInt(dbPortEnvKey, 5432)
	dbUser := getEnvironment(dbUserEnvKey, "pguser")
	dbPassword := getEnvironment(dbPasswordEnvKey, "pgpwd")
	dbName := getEnvironment(dbNameEnvKey, "urlshortner")
	dbSSLMode := getEnvironment(dbSSLModeEnvKey, "disable")
	maxOpenConnections := getEnvironmentAsInt(maxOpenConnectionsEnvKey, 25)
	maxIdleConnections := getEnvironmentAsInt(maxIdleConnectionsEnvKey, 10)
	connectionMaxLifetime := getEnvironmentAsInt(connectionMaxLifetimeEnvKey, 30)
	connectionMaxIdleTime := getEnvironmentAsInt(connectionMaxIdleEnvKey, 10)
	applicationPort := getEnvironmentAsInt(applicationPortEnvKey, 8080)

	return Config{
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		dbSSLMode,
		maxOpenConnections,
		maxIdleConnections,
		connectionMaxLifetime,
		connectionMaxIdleTime,
		applicationPort,
	}
}

func getEnvironment(key, defaultValue string) string {
	if envValue := os.Getenv(key); envValue != "" {
		return envValue
	}
	return defaultValue
}

func getEnvironmentAsInt(key string, defaultValue int) int {
	if envValue := os.Getenv(key); envValue != "" {
		if intValue, err := strconv.Atoi(envValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}
