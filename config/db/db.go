package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config is the db config
type Config struct {
	Dialect            string
	Username, Password string
	DbName, DbHost     string
	DbURL              string
}

var config map[string]Config

// GetConfig returns DB config based on APP_ENV variable
func GetConfig(appEnv string) Config {
	return config[appEnv]
}

func init() {
	err := godotenv.Load(os.ExpandEnv("$GOPATH/src/github.com/dikaeinstein/go-rest-api/.env"))
	if err != nil {
		log.Println(err)
	}
	config = map[string]Config{
		"development": Config{
			Dialect:  os.Getenv("DB_TYPE"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			DbName:   os.Getenv("DB_DEV"),
			DbHost:   os.Getenv("DB_HOST"),
		},
		"test": Config{
			Dialect:  os.Getenv("DB_TYPE"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			DbName:   os.Getenv("DB_TEST"),
			DbHost:   os.Getenv("DB_HOST"),
		},
		"production": Config{
			Dialect: os.Getenv("DB_TYPE"),
			DbURL:   os.Getenv("DATABASE_URL"),
		},
	}
}
