package config

import (
	"log"
	"net"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// DBConfig holds the database config options
type DBConfig struct {
	Dialect            string
	Username, Password string
	DbName, DbHost     string
	DbURL              string
	Logging            bool
}

var dbConfig map[string]DBConfig

// GetConfig returns DB config based on APP_ENV variable
func getDbConfig(appEnv string) DBConfig {
	dbURL := dbConfig[appEnv].DbURL
	if dbURL == "" && appEnv != "production" {
		return DBConfig{
			Dialect:  getEnv("DB_TYPE", "postgres"),
			Username: getEnv("DB_USER", ""),
			Password: getEnv("DB_PASS", ""),
			DbName:   getEnv("DB_DEV", ""),
			DbHost:   getEnv("DB_HOST", "localhost"),
			Logging:  getEnvAsBool("DB_LOGGING", false),
		}
	}
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	password, _ := parsedURL.User.Password()
	host, _, _ := net.SplitHostPort(parsedURL.Host)
	return DBConfig{
		Username: parsedURL.User.Username(),
		Password: password,
		DbName:   parsedURL.Path[1:],
		DbHost:   host,
		Logging:  getEnvAsBool("DB_LOGGING", false),
	}
}

func init() {
	err := godotenv.Load(os.ExpandEnv("$GOPATH/src/github.com/dikaeinstein/go-rest-api/.env"))
	if err != nil {
		log.Println(err)
	}
}
