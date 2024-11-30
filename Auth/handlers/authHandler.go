package handlers

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	tokens "AuthServ/Tokens"
	"AuthServ/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func AuthHandler(ctx *gin.Context) {
	GUID := ctx.Param("guid")            // получаем GUID из параметра запроса
	var sessionUser models.ClientSession // структура для сохранения RefreshToken в БД
	var createdTokens models.Tokens      // структура для отправки токенов
	AdressIp := "127.0.0.1"              // данные IP-адреса для примера
	// AdressIp := ctx.ClientIP() 		 // получаем IP-адрес клиента, если бы он был

	// создаем AcessToken
	createdAcessToken, err := tokens.GenerateAcessToken(GUID, AdressIp)
	if err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to generate access token"})
		log.Panicln(err)
		return
	}
	
	// создаем RefreshToken
	createdRefreshToken, err := tokens.GenerateRefreshToken()
	if err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to generate refresh token"})
		log.Panicln(err)
		return
	}

	// хешируем RefreshToken, так как -> хранится в базе исключительно в виде bcrypt хеша
	hashedRefreshToken, err := utils.Hashing(createdRefreshToken)
	if err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to hash the refresh token"})
		log.Panicln(err)
		return
	}

	// заполняем данные для записи в БД
	sessionUser.RefreshToken = hashedRefreshToken // в БД отправляется bcrypt хеш
	sessionUser.SessionGUID = GUID
	sessionUser.SessionIP = AdressIp

	// сохраняем данные
	if res := database.DB.Create(&sessionUser); res.Error != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to save to the database"})
		log.Panicln(res.Error)
		return
	}
	log.Println("Запись в БД успешно создана")

	// заполняем данные для ответа 
	createdTokens.AccessToken = createdAcessToken
	createdTokens.RefreshToken = createdRefreshToken // в ответе токен отправляется по формату base64

	ctx.JSON(200, models.ResponceData{Message: "tokens successfully issued", Data: createdTokens})
}
