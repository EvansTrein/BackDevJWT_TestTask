package handlers

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Удаление пользователя
// @Description Удаление пользователя с указанным GUID и его сессии
// @Tags user
// @Accept json
// @Produce json
// @Param guid path string true "Уникальный идентификатор пользователя" Format(7c5e66cf-57ba-4871-9186-74ff5ab1e1f1)
// @Success 200 {object} models.ResponceMessage
// @Failure 400 {object} models.ErrResponce
// @Failure 404 {object} models.ErrResponce
// @Failure 500 {object} models.ErrResponce
// @Security accessToken
// @Security refreshRefresh
// @Router /user/{guid}/del [delete]
func UserDelHandler(ctx *gin.Context) {
	GUID := ctx.Param("guid")            // получаем GUID из параметра запроса
	var acriveUser models.User           // пользователь, которого будем удалять
	var sessionUser models.ClientSession // сессия пользователя, которую будем удалять

	// ищем нужного пользователя
	if resFind := database.DB.Where("guid = ?", GUID).First(&acriveUser); resFind.Error != nil {
		if errors.Is(resFind.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(404, models.ErrResponce{ErrMessage: "user not found"})
			return
		}
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to find user"})
		return
	}

	// удаление пользователя
	if resDel := database.DB.Unscoped().Delete(&acriveUser); resDel.Error != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to delete user"})
		return
	}
	log.Println("пользователь успешно удален")

	// ищем сессию
	if resFind := database.DB.Where("session_guid = ?", GUID).First(&sessionUser); resFind.Error != nil {
		if errors.Is(resFind.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(404, models.ResponceMessage{Message: "user not found"})
			return
		}
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to find session"})
		return
	}

	// удаляем сессию
	if resDel := database.DB.Unscoped().Delete(&sessionUser); resDel.Error != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to delete session"})
		return
	}
	log.Println("сессия пользователя успешно удалена")

	ctx.JSON(200, models.ResponceMessage{Message: "user and session successfully deleted"})
}
