package handlers

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Удаление сессии
// @Description Удаление сессии с указанным GUID, если нужно удалить вручную
// @Tags auth
// @Accept json
// @Produce json
// @Param guid path string true "Уникальный идентификатор сессии" Format(7c5e66cf-57ba-4871-9186-74ff5ab1e1f1)
// @Success 200 {object} models.ResponceMessage
// @Failure 400 {object} models.ErrResponce
// @Failure 404 {object} models.ErrResponce
// @Failure 500 {object} models.ErrResponce
// @Router /auth/delSession/{guid} [delete]
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

	log.Println("сессия успешно удалена")
	ctx.JSON(200, models.ResponceMessage{Message: "session successfully deleted"})
}
