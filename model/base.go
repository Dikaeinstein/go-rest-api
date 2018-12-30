package model

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"

	"github.com/dikaeinstein/go-rest-api/config/db"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Register the postgres db driver
)

var d *gorm.DB // Database

func init() {
	appEnv := os.Getenv("APP_ENV")
	dbURI := parseDbConfig(appEnv)
	dialect := db.GetConfig(appEnv).Dialect
	connectDB(dialect, dbURI)
	// d.LogMode(true)
}

func connectDB(dialect, dbURI string) {
	fmt.Println(dbURI)

	conn, err := gorm.Open(dialect, dbURI)
	if err != nil {
		log.Println(err)
	}

	d = conn
	// Database migration
	d.Debug().AutoMigrate(&Account{}, &Contact{})
}

func parseDbConfig(appEnv string) string {
	var username, password, dbName, dbHost string

	if appEnv == "production" || appEnv == "prod" {
		parsedURL, err := url.Parse(db.GetConfig(appEnv).DbURL)
		if err != nil {
			log.Println(err)
		}
		username = parsedURL.User.Username()
		password, _ = parsedURL.User.Password()
		dbName = parsedURL.Path[1:]
		dbHost, _, _ = net.SplitHostPort(parsedURL.Host)
	} else {
		c := db.GetConfig(appEnv)
		username = c.Username
		password = c.Password
		dbName = c.DbName
		dbHost = c.DbHost
	}

	return fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		dbHost, username, dbName, password)
}

// GetDB returns a handle to the DB object
func GetDB() *gorm.DB {
	return d
}
