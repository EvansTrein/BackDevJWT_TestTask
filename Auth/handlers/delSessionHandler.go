package handlers

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DelSessionHandler(ctx *gin.Context) {
	GUID := ctx.Param("guid")            // получаем GUID из параметра запроса
	var sessionUser models.ClientSession // структура сессии, которую будем удалять

	if resFind := database.DB.Where("session_guid = ?", GUID).First(&sessionUser); resFind.Error != nil {
		if errors.Is(resFind.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(404, models.ErrResponce{ErrMessage: "session not found"})
			return
		}
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to find session"})
		return
	}

	if resDel := database.DB.Unscoped().Delete(&sessionUser); resDel.Error != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to delete session"})
		return
	}

	ctx.JSON(200, models.ResponceMessage{Message: "session successfully deleted"})
}
