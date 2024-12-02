package handlers

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	"AuthServ/utils"
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserUpdateHandler(ctx *gin.Context) {
	var activeUser models.User
	var newDataUser models.UserDataUpd
	GUID := ctx.Param("guid")

	// получаем данные из тела запроса
	err := ctx.BindJSON(&newDataUser)
	if err != nil {
		ctx.JSON(400, models.ErrResponce{ErrMessage: "incorrect data in body"})
		return
	}

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

	// проверяем отличается ли email
	if newDataUser.Email != activeUser.EmailUser {
		activeUser.EmailUser = newDataUser.Email
	}

	// хешируем пароль
	hashedPass, err := utils.Hashing(newDataUser.Pass)
	if err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "password hashing failed"})
		return
	}

	activeUser.PassUserHash = hashedPass // сохраняем хешированный пароль

	// обновляем запись в БД
	if resUpd := database.DB.Save(&activeUser); resUpd.Error != nil {
		if strings.Contains(resUpd.Error.Error(), "23505") { // это проверка на дубликаты
			ctx.JSON(400, models.ErrResponce{ErrMessage: "this email already exists"})
			return
		}
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to update user"})
		log.Printf("ERROR - не удалось сохранить нового пользователя\nОшибка: %s", resUpd.Error)
		return
	}

	log.Println("данные пользователя успешно обновлены")
	ctx.JSON(200, models.ResponceMessage{Message: "user data successfully updated"})
}
