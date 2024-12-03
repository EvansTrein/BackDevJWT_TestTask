package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {

	uri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		"db_AuthServ", "evans", "evans", "postgres", "4010", "disable")

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return err
	} else {
		DB = db
		return nil
	}
}
