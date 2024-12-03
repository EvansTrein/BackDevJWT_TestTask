package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {

	uri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_NAME"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USE_SSL"))

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return err
	} else {
		DB = db
		return nil
	}
}
