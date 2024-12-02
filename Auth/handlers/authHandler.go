package handlers

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	tokens "AuthServ/Tokens"
	"AuthServ/utils"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthHandler(ctx *gin.Context) {
	GUID := ctx.Param("guid")            // получаем GUID из параметра запроса
	var sessionUser models.ClientSession // структура для сохранения RefreshToken в БД
	var createdTokens models.Tokens      // структура для отправки токенов
	var activeUser models.User           // структура для пользователя, который ранее зарегистрировался
	var authData models.AuthData         // структура для данных, которые передали для входа
	AdressIp := "127.0.0.1"              // данные IP-адреса для примера
	// AdressIp := ctx.ClientIP() 		 // получаем IP-адрес клиента, если бы он был

	// ищем пользователя
	if findUser := database.DB.Where("guid = ?", GUID).First(&activeUser); findUser.Error != nil {
		if errors.Is(findUser.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(404, models.ErrResponce{ErrMessage: "not found user"})
			return
		} else {
			ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to search the user database"})
			return
		}
	}

	// парсим данные из тела запроса
	if err := ctx.BindJSON(&authData); err != nil {
		ctx.JSON(400, models.ErrResponce{ErrMessage: "invalid request body"})
		return
	}

	// проверяем, что GUID пришедший в Param совпадает с GUID, который нашли у пользователя из БД, прост. на всякий случай
	if activeUser.GUID != GUID {
		ctx.JSON(400, models.ErrResponce{ErrMessage: "guid in the parameters does not match the guid of the user found in the database"})
		return
	}

	// проверяем email
	if activeUser.EmailUser != authData.Email {
		ctx.JSON(400, models.ErrResponce{ErrMessage: "invalid email"})
		return
	}

	// сверяем пришедший пароль с хешем пароля, который хранится в БД
	if isPass := utils.CheckHashing(authData.Pass, activeUser.PassUserHash); !isPass {
		ctx.JSON(400, models.ErrResponce{ErrMessage: "incorrect password"})
		return
	}

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
	sessionUser.MaxSessionDuration = time.Duration(time.Duration(3600 * time.Second)) // устанавливаем время жизни токена, тут 1 час
	log.Println(sessionUser.MaxSessionDuration)

	// сохраняем данные
	if res := database.DB.Create(&sessionUser); res.Error != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to save to the database"})
		log.Panicln(res.Error)
		return
	}
	log.Println("Запись сессии в БД успешно создана")

	// заполняем данные для ответа
	createdTokens.AccessToken = createdAcessToken
	createdTokens.RefreshToken = createdRefreshToken // в ответе токен отправляется по формату base64

	ctx.JSON(200, models.ResponceData{Message: "tokens successfully issued", Data: createdTokens})
}
