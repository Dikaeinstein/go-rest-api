package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds the application configuration options
type Config struct {
	Db   DBConfig
	Port int
}

// New returns a new config struct
func New() *Config {
	if appEnv, exists := os.LookupEnv("APP_ENV"); exists {
		return &Config{Db: getDbConfig(appEnv)}
	}
	return &Config{
		Db:   getDbConfig("development"),
		Port: getEnvAsInt("PORT", 4000),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper to read an environment variable into a bool or return a default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

// Simple helper function to read an environment variable into an integer
// or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// Helper to read an environment variable into a string slice or return a default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")
	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)
	return val
}
