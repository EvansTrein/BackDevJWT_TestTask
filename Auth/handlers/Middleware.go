package handlers

import (
	models "AuthServ/Models"
	tokens "AuthServ/Tokens"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Middleware(ctx *gin.Context) {
	// получаем AccessToken из заголовка Authorization и обрезаем префикс Bearer
	incomingAccessToken := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
	// получаем RefreshToken из заголовка RefreshToken
	incomingRefreshToken := ctx.GetHeader("RefreshToken")

	if incomingAccessToken == "" {
		ctx.JSON(401, models.ErrResponce{ErrMessage: "unauthorized user"})
		return
		// ctx.Redirect(303, "/login")
	}

	// получаем из AccessToken сам токен и проверяем его
	oldAccessToken, err := tokens.ValidateAccessToken(incomingAccessToken)
	if err != nil {
		ctx.JSON(400, models.ErrResponce{ErrMessage: err.Error()})
		return
	}

	// проверка валидности AccessToken
	if oldAccessToken.Valid {
		ctx.JSON(200, models.ResponceMessage{Message: "AccessToken has been successfully"})
		log.Println("AccessToken успешно прошел")
		ctx.Next()
		return
	}

	// Сюда доходим, если AccessToken истек, вызвыаем Refresh опеарцию

	// Создаем новый HTTP-запрос
	req, err := http.NewRequestWithContext(ctx.Request.Context(), "POST", "http://localhost:4000/auth/refresh", ctx.Request.Body)
	if err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to create a request to update tokens"})
		ctx.Abort()
		return
	}

	// Добавляем необходимые заголовки
	req.Header.Set("Authorization", incomingAccessToken)
	req.Header.Set("RefreshToken", incomingRefreshToken)

	// Отправляем запрос на внешний сервис
	client := &http.Client{}
	respRefreshOperation, err := client.Do(req)
	if err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to send a token update request"})
		ctx.Abort()
		return
	}
	defer respRefreshOperation.Body.Close()

	// читаем тело запроса
	body, err := io.ReadAll(respRefreshOperation.Body)
	if err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to read response body"})
		ctx.Abort()
		return
	}

	// проверяем результат Refresh операции
	switch {
	case respRefreshOperation.StatusCode == 500:
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to update tokens"})
		ctx.Abort()
		return
	case respRefreshOperation.StatusCode != 200:
		ctx.JSON(401, models.ErrResponce{ErrMessage: "RefreshToken expired, need to re-authorize"})
		ctx.Abort()
		// ctx.Redirect(303, "/login")
		return
	}

	// приводим данные из тела запроса к ключ-значение
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, models.ErrResponce{ErrMessage: "error when converting answer body to map"})
		ctx.Abort()
		return
	}

	// получаем новые токены
	newTokens, ok := data["Data"]
	if !ok {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "no new tokens came from the Refresh operation"})
		ctx.Abort()
		return
	}

	log.Println("через Middleware были созданы новое токены")
	ctx.JSON(200, models.ResponceData{Message: "tokens have been updated", Data: newTokens})
	ctx.Next()
}
