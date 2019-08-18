package model

import (
	"fmt"
	"log"

	"github.com/dikaeinstein/go-rest-api/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Register the postgres db driver
)

var d *gorm.DB // Database

func init() {
	config := config.New()
	connectDB(config.Db)
	d.LogMode(config.Db.Logging)
}

func New(db *gorm.DB) {

}

func connectDB(dc config.DBConfig) {
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		dc.DbHost, dc.Username, dc.DbName, dc.Password)

	conn, err := gorm.Open(dc.Dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	}

	d = conn
	// Database migration
	d.Debug().AutoMigrate(&Account{}, &Contact{})
}

// GetDB returns a handle to the DB object
func GetDB() *gorm.DB {
	return d
}
