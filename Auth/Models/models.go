package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GUID         string `gorm:"not null;unique"`
	EmailUser    string `gorm:"not null;unique"`
	PassUserHash string `gorm:"pass" json:"-"`
}

type UserData struct {
	GUID      string `json:"guid" binding:"required" example:"7c5e66cf-57ba-4871-9186-74ff5ab1e1f1"`
	EmailUser string `json:"email" binding:"required,email" example:"user1@mail.com"`
	PassUser  string `json:"password" binding:"required,min=6" example:"123456"`
}

type UserDataUpd struct {
	Email string `json:"email" binding:"required,email" example:"user1@mail.com"`
	Pass  string `json:"password" binding:"required,min=6" example:"123456"`
}

type AuthData struct {
	Email string `json:"email" binding:"required,email" example:"user1@mail.com"`
	Pass  string `json:"password" binding:"required,min=6" example:"123456"`
}

type ClientSession struct {
	gorm.Model
	RefreshToken       string        `gorm:"not null;unique" json:"-"`
	SessionGUID        string        `gorm:"not null;unique" json:"-"`
	SessionIP          string        `gorm:"not null" json:"-"`
	RefreshTokenID     uint          `gorm:"not null" json:"-"`
	MaxSessionDuration time.Duration `gorm:"not null" json:"-"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRJUCI6IjEyNy4wLjAuMSIsImV4cCI6MTczMzE4MzY5MSwicmVmcmVzaFRva2VuSUQiOjExfQ.LyuwUe7IPSG2_aPdT59Ms2_xmDPa9-ymhGsuuJ_uwi5wzxfjoHerNSTpJLf2ZQUXGNjDp3BHgs2jXw4ehLLjuQ"`
	RefreshToken string `json:"refreshToken" example:"3r65EyQIo/NsGR3TE1/Y7GIuD+jm1diGf+zZ4DoXwhg="`
}

type ErrResponce struct {
	ErrMessage string `example:"error message"`
}

type ResponceData struct {
	Message string      `example:"info message"`
	Data    interface{} `example:"responce data"`
}

type ResponceMessage struct {
	Message string `example:"info message"`
}

type RespGetMidlExample struct {
	MiddlewareStatus string `example:"тут могут быть сообщения от Middleware или ошибки или новые токены"` 
	Guid             string `example:"3c43e84d-fc44-4895-bc72-2a7924417b80"`
	Email            string `example:"user1@mail.com"`
}

type RespPutMidlExample struct {
	MiddlewareStatus string `example:"тут могут быть сообщения от Middleware или новые токены"`
	ResultPutMess string `example:"info message"`
}
