package database

import (
	"github.com/jinzhu/gorm"
)

// User database structure
type User struct {
	gorm.Model
	ID             int64
	Email          string
	Password       string
	SessionToken   string
	FailedAttempts int
}
