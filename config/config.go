package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
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
}

const dbHostEnvKey = "DB_HOST"
const dbPortEnvKey = "DB_PORT"
const dbUserEnvKey = "DB_USER"
const dbPasswordEnvKey = "DB_PASSWORD"
const dbNameEnvKey = "DB_NAME"
const dbSSLModeEnvKey = "DB_SSL_MODE"
const maxOpenConnectionsEnvKey = "DB_MAX_OPEN_CONNECTIONS"
const maxIdleConnectionsEnvKey = "DB_MAX_IDLE_CONNECTIONS"
const connectionMaxLifetimeEnvKey = "DB_CONN_MAX_LIFETIME_MIN"
const connectionMaxIdleEnvKey = "DB_CONN_MAX_IDLE_TIME_MIN"

func LoadDBConfig() DBConfig {
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

	return DBConfig{
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		dbSSLMode,
		maxOpenConnections,
		maxIdleConnections,
		connectionMaxLifetime,
		connectionMaxIdleTime}
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
