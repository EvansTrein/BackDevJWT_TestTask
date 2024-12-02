package handlers

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	tokens "AuthServ/Tokens"
	"AuthServ/utils"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func delSession(guid string, table *models.ClientSession) error {

	// удаляем старую запись с RefreshToken из БД, удаляем без возможности восстановления
	// получается, что использовать старый RefreshToken повторно не получится, так как он удаляется
	res := database.DB.Where("session_guid = ?", guid).Unscoped().Delete(table)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func AuthRefreshHandler(ctx *gin.Context) {
	var activeSission models.ClientSession    // переменная для хранения активной сессии, в которой находится RefreshToken
	var newActiveSission models.ClientSession // переменная для новой сессии, которая будет создана при обнолении токенов
	var newTokens models.Tokens               // переменная для новых токенов
	incomingIP := "127.1.1.1"                 // переменная для ip адреса, с которого был зовершен запрос
	emailUser := "evanstrein@icloud.com"      // для отправки warning на почту, по хорошему, получаем из базы данных
	// incomingIP := ctx.ClientIP()			  // получаем IP-адрес клиента, если бы он был

	// получаем AccessToken из заголовка Authorization и обрезаем префикс Bearer
	incomingAccessToken := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
	// получаем RefreshToken из заголовка RefreshToken
	incomingRefreshToken := ctx.GetHeader("RefreshToken")

	// проверка, что нужные данные были переданы
	if incomingAccessToken == "" || incomingRefreshToken == "" {
		ctx.JSON(400, models.ErrResponce{ErrMessage: "accessToken or refreshToken were not passed in the request headers"})
		return
	}

	// получаем из AccessToken сам токен и проверяем его
	oldAccessToken, err := tokens.ValidateAccessToken(incomingAccessToken)
	if err != nil {
		ctx.JSON(400, models.ErrResponce{ErrMessage: err.Error()})
		return
	} else if oldAccessToken.Valid {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "was triggered to update tokens with still valid AccessToken"})
		log.Println("ERROR - получен действительный AccessToken для обновления")
		return
	}

	// вытаскиваем нагрузку из AccessToken
	paylodAccessToken := oldAccessToken.Claims.(jwt.MapClaims)
	guidSession, ok := paylodAccessToken["guid"]
	if !ok {
		ctx.JSON(400, models.ErrResponce{ErrMessage: "failed to get a guid from oldAccessToken"})
		return
	}
	guidSessionStr, ok := guidSession.(string) // из нагрузки мы получили данные как interface, а нам нужна string
	if !ok {
		ctx.JSON(400, models.ErrResponce{ErrMessage: "failed to convert guid to string type"})
		return
	}

	// ищем по GUID из AccessToken запись в БД с RefreshToken
	if res := database.DB.Where("session_guid = ?", guidSessionStr).First(&activeSission); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(404, models.ErrResponce{ErrMessage: "failed to find a session with such guid"})
			return
		} else {
			ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to search the active sessions database"})
			return
		}
	}

	// проверяем пришедший в запросе RefreshToken(пришла строка) с хешем из записи в БД с RefreshToken
	// при создании, в нагрузке AccessToken был GUID и при создании RefreshToken - в запись добавляли GUID
	// получается, что GUID - общий уникальный идентификатор, который обоюдно связывает AccessToken и RefreshToken
	// значит, если по этому GUID найдется ДРУГОЙ хеш RefreshToken - токены не были созданы вместе
	if isValidHash := utils.CheckHashing(incomingRefreshToken, activeSission.RefreshToken); !isValidHash {
		ctx.JSON(400, models.ErrResponce{ErrMessage: "this refreshToken cannot be used to issue new tokens"})
		// ctx.Redirect(303, "/login") // нужно заново войти, перенаправляем
		return
	}

	// еще одна проверка, для перестраховки, если выше хеш из RefreshToken, полученный по GUID из AccessToken
	// прошел проверку, но при этом, GUID были разные - это очень плохо!
	if guidSessionStr != activeSission.SessionGUID {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "the GUID in AccessToken is different from the GUID in RefreshToken"})
		// если попали сюда, значит произошло то, что считалось невозможным
		// в этом случае надо какие-то действия предпринять, сообщить кому-то например
		return
	}

	// получаем прошедшее время с момента создания токена
	sinceRefreshTokenCreated := time.Since(activeSission.CreatedAt)

	// сравниваем прошедшее время с максимальным временем жизни токена
	if sinceRefreshTokenCreated > activeSission.MaxSessionDuration {

		// удаляем сессию из БД, если время жизни RefreshToken истекло
		if err := delSession(guidSessionStr, &newActiveSission); err != nil {
			ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to delete an old session in the database"})
			log.Panicln(err)
			return
		}

		log.Println("Сессия у которой вышло время успешно удалена")
		ctx.JSON(401, models.ErrResponce{ErrMessage: "the refreshToken's lifetime has expired and it has been deleted"})
		// ctx.Redirect(303, "/login") // если RefreshToken истек, то нужно перенаправить на страницу входа
		return
	}

	// проверка ip-адресов, если новый -> посылаем email warning на почту
	sendWarning := make(chan string, 1)
	if incomingIP != activeSission.SessionIP {
		log.Println("новый ip, отправка письма!")
		go func() {
			resultSend, errSend := utils.SendEmailWarning(emailUser, incomingIP)
			if resultSend != "" {
				sendWarning <- resultSend
			} else {
				sendWarning <- errSend
			}
			close(sendWarning)
		}()
	} else {
		sendWarning <- "the ip address has not been changed"
	}

	// создаем новый AcessToken
	createdNewAcessToken, err := tokens.GenerateAcessToken(guidSessionStr, incomingIP)
	if err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to generate access token"})
		log.Panicln(err)
		return
	}

	// создаем новый RefreshToken
	createdNewRefreshToken, err := tokens.GenerateRefreshToken()
	if err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to generate refresh token"})
		log.Panicln(err)
		return
	}

	// хешируем новый RefreshToken, так как -> хранится в базе исключительно в виде bcrypt хеша
	hashedNewRefreshToken, err := utils.Hashing(createdNewRefreshToken)
	if err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to hash the refresh token"})
		log.Panicln(err)
		return
	}

	// готовим данные для записи новой сессии
	newActiveSission.RefreshToken = hashedNewRefreshToken // в БД отправляется bcrypt хеш
	newActiveSission.SessionGUID = guidSessionStr
	newActiveSission.SessionIP = incomingIP
	newActiveSission.MaxSessionDuration = time.Duration(time.Duration(3600 * time.Second)) // устанавливаем время жизни токена, тут 1 час

	// проверяем результат отправки email warning, код не пойдет дальше, пока не будут получены данные из канала
	resSendEmail := <-sendWarning
	switch resSendEmail {
	case "email warning was successfully sent":
		log.Println("email warning был успешно отправлен")
	case "the ip address has not been changed":
		log.Println("ip-адрес не менялся")
	default:
		log.Printf("ERROR - email warning не был отправлен\nОшибка: %s", resSendEmail)
	}

	// удаляем страую сессию
	if err := delSession(guidSessionStr, &newActiveSission); err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to delete an old session in the database"})
		log.Panicln(err)
		return
	}
	log.Println("Запись старой сессии в БД успешно удалена")

	// создаем новую сессию в БД, с новым RefreshToken и обновленным ip-адресом (если он действительно новый)
	if res := database.DB.Create(&newActiveSission); res.Error != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to save a new session in the database"})
		log.Panicln(res.Error)
		return
	}
	log.Println("Запись обновленной сессии в БД успешно создана")

	// записываем новые токены для ответа
	newTokens.AccessToken = createdNewAcessToken
	newTokens.RefreshToken = createdNewRefreshToken // в ответе токен отправляется по формату base64

	ctx.JSON(200, models.ResponceData{Message: "token update was successful ", Data: newTokens})
	log.Println("обновление токенов прошло успешно ")
}
