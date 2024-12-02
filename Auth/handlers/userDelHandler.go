package handlers

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserDelHandler(ctx *gin.Context) {
	GUID := ctx.Param("guid")  // получаем GUID из параметра запроса
	var acriveUser models.User // пользователь, которого будем удалять

	if resFind := database.DB.Where("guid = ?", GUID).First(&acriveUser); resFind.Error != nil {
		if errors.Is(resFind.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(404, models.ErrResponce{ErrMessage: "user not found"})
			return
		}
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to find user"})
		return
	}

	if resDel := database.DB.Unscoped().Delete(&acriveUser); resDel.Error != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to delete user"})
		return
	}

	ctx.JSON(200, models.ResponceMessage{Message: "user successfully deleted"})
}
