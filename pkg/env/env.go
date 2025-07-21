package env

import (
	"os"
	"strconv"
)

// Postgres connection variables
var PostgresUser = GetEnvVar("etzba", "POSTGRES_USER")
var PostgresPass = GetEnvVar("Pass1234", "POSTGRES_PASSWORD")
var PostgresDB = GetEnvVar("etzba", "POSTGRES_DB")
var PostgresHost = GetEnvVar("localhost", "POSTGRES_HOST")
var PostgresPort = GetEnvInt(5432, "POSTGRES_PORT")

func GetEnvVar(defaultValue, key string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func GetEnvInt(defaultValue int, key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value == 0 {
		return defaultValue
	}
	return value
}
