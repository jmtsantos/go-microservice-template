package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres dialect
)

// Database type
type Database struct {
	db *gorm.DB
}

// New returns a new database handler
func New(domain, port string) *Database {
	var (
		database Database
		err      error
	)

	// Connect to postgres
	if database.db, err = gorm.Open("postgres", connectString()); err != nil {
		panic(err)
	}
	defer database.db.Close()

	// Auto migrate structs
	database.db.AutoMigrate(
		&User{},
	)

	return &database
}

func connectString() string {
	return fmt.Sprintf("host=postgres user=postgres dbname=go-template password=changeme sslmode=disable")
}
