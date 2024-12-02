package handlers

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserHandler(ctx *gin.Context) {
	var activeUser models.User // переменная для поиска пользователя
	GUID := ctx.Param("guid")  // получаем GUID из параметра запроса

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

	// формируем ответ
	resp := models.ResponceData{
		Message: "user data successfully found",
		Data: struct {
			GUID  string `json:"guid"`
			Email string `json:"email"`
		}{
			GUID:  activeUser.GUID,
			Email: activeUser.EmailUser,
		},
	}

	ctx.JSON(200, resp)
	log.Println("данные пользователя успешно получены")
}
