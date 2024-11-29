package server

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	"log"
)

func InitServer() {

	errDB := database.InitDatabase()
	if errDB != nil {
		log.Fatal("Ошибка подключения к базе данных: ", errDB)
	} else {
		log.Println("Успешное подключение к базе данных")
		// Миграции БД
		database.DB.AutoMigrate(&models.User{})
		database.DB.AutoMigrate(&models.ClientSession{})
	}

}

func StartServer() {
	InitRotes()
}
