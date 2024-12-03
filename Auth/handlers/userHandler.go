package handlers

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Получение данных пользователя
// @Description Получение данных пользователя с указанным GUID
// @Tags user
// @Accept json
// @Produce json
// @Param guid path string true "Уникальный идентификатор пользователя" Format(7c5e66cf-57ba-4871-9186-74ff5ab1e1f1)
// @Success 200 {object} models.RespGetMidlExample
// @Failure 400 {object} models.ErrResponce
// @Failure 404 {object} models.ErrResponce
// @Failure 500 {object} models.ErrResponce
// @Security accessToken
// @Security refreshRefresh
// @Router /user/{guid} [get]
func UserHandler(ctx *gin.Context) {
	var activeUser models.User   // переменная для поиска пользователя
	var resp models.ResponceData // переменная для отправки ответа
	GUID := ctx.Param("guid")    // получаем GUID из параметра запроса

	// поиск в базе
	if searchRes := database.DB.Where("guid = ?", GUID).First(&activeUser); searchRes.Error != nil {
		if errors.Is(searchRes.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(404, models.ErrResponce{ErrMessage: "not found user"})
			return
		} else {
			ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to search the user database"})
			return
		}
	}

	// с обработчика Middleware могут прийти такие ключи, там будет результат работы Middleware
	MiddlewareInfo, okMidlInfo := ctx.Get("MessageMiddleware")
	MiddlewareData, okMidlData := ctx.Get("DataMiddleware") // если есть, значит пришли новые токены

	// добавляем к ответу этого обработчика инфо от Middleware, если оно было и не было новых токенов
	if okMidlInfo {
		info, ok := MiddlewareInfo.(models.ResponceMessage)
		if !ok {
			ctx.JSON(500, models.ErrResponce{ErrMessage: "error converting the value from MessageMiddleware key to models.ResponceMessage"})
			log.Panicln("ошибка приведения зачения из ключа MessageMiddleware к models.ResponceMessage")
			return
		}

		// формируем ответ
		resp.Message = "user data successfully found"
		resp.Data = struct {
			MiddlewareInfo string `json:"MiddlewareInfo"`
			GUID           string `json:"guid"`
			Email          string `json:"email"`
		}{
			MiddlewareInfo: info.Message,
			GUID:           activeUser.GUID,
			Email:          activeUser.EmailUser,
		}
	}

	// добавляем к ответу этого обработчика инфо от Middlewar и новые токены, если они были обновлены
	if okMidlData {
		data, ok := MiddlewareData.(models.ResponceData)
		if !ok {
			ctx.JSON(500, models.ErrResponce{ErrMessage: "error converting values from the DataMiddleware key to models.ResponceData"})
			log.Panicln("ошибка приведения зачения из ключа DataMiddleware к models.ResponceData")
			return
		}

		// формируем ответ
		resp.Message = "user data successfully found"
		resp.Data = struct {
			MiddlewareData models.ResponceData `json:"MiddlewareData"`
			GUID           string              `json:"guid"`
			Email          string              `json:"email"`
		}{
			MiddlewareData: data,
			GUID:           activeUser.GUID,
			Email:          activeUser.EmailUser,
		}
	}

	ctx.JSON(200, resp)
	log.Println("данные пользователя успешно получены")
}
