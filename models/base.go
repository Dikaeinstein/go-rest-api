package models

import (
	"fmt"
	"net"
	"os"

	"net/url"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Register the postgres db driver
	"github.com/joho/godotenv"
)

var db *gorm.DB // Database

func init() {
	err := godotenv.Load() // Load .env file
	if err != nil {
		fmt.Print(err)
	}

	parsedURL, err := url.Parse(os.Getenv("DATABASE_URL"))

	dialect := parsedURL.Scheme
	username := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()
	dbName := parsedURL.Path[1:]
	dbHost, _, _ := net.SplitHostPort(parsedURL.Host)

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	fmt.Println(dbURI)

	conn, err := gorm.Open(dialect, dbURI)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &Contact{}) // Database migration
}

// GetDB returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
