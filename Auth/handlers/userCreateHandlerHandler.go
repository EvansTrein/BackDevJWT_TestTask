package handlers

import (
	database "AuthServ/Database"
	models "AuthServ/Models"
	"AuthServ/utils"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Создание нового пользователя
// @Description Создание нового пользователя с предоставленными данными
// @Tags user
// @Accept json
// @Produce json
// @Param body body models.UserData true "Данные пользователя"
// @Success 201 {object} models.ResponceMessage
// @Failure 400 {object} models.ErrResponce
// @Failure 500 {object} models.ErrResponce
// @Router /userCreate [post]
func UserCreateHandler(ctx *gin.Context) {
	var newUser models.User      // переменная для запси в БД user
	var userData models.UserData // переменная для данных из запроса

	// получаем данные из тела запроса
	err := ctx.BindJSON(&userData)
	if err != nil {
		ctx.JSON(400, models.ErrResponce{ErrMessage: "incorrect data in body"})
		return
	}

	// проверяем формат GUID
	if isVaildGUID := utils.IsGUID(userData.GUID); !isVaildGUID {
		ctx.JSON(400, models.ErrResponce{ErrMessage: "incorrect GUID format"})
		return
	}

	// хешируем пароль
	hashedPass, err := utils.Hashing(userData.PassUser)
	if err != nil {
		ctx.JSON(500, models.ErrResponce{ErrMessage: "password hashing failed"})
		return
	}

	// готовим данные для записи в БД
	newUser.GUID = userData.GUID
	newUser.EmailUser = userData.EmailUser
	newUser.PassUserHash = hashedPass

	// создаем запись в БД
	if resSave := database.DB.Create(&newUser); resSave.Error != nil {
		if strings.Contains(resSave.Error.Error(), "23505") { // это проверка на дубликаты, в БД guid и emai должны быть уникальны
			ctx.JSON(400, models.ErrResponce{ErrMessage: "user with this guid or email already exists"})
			return
		}
		ctx.JSON(500, models.ErrResponce{ErrMessage: "failed to save a new user"})
		log.Printf("ERROR - не удалось сохранить нового пользователя\nОшибка: %s", resSave.Error)
		return
	}

	ctx.JSON(201, models.ResponceMessage{Message: "new user successfully saved"})
	log.Println("новый пользователь успешно сохранен")
}
