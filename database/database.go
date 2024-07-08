package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Connect() error {
	var err error
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_DATABASE")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos", host, username, password, databaseName, port)
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	return nil
}
