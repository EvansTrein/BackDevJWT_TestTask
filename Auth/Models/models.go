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
	GUID      string `json:"-"`
	EmailUser string `json:"email" binding:"required,email"`
	PassUser  string `json:"password" binding:"required,min=6"`
}

type ClientSession struct {
	gorm.Model
	RefreshToken       string        `gorm:"not null;unique" json:"-"`
	SessionGUID        string        `gorm:"not null;unique" json:"-"`
	SessionIP          string        `gorm:"not null" json:"-"`
	MaxSessionDuration time.Duration `gorm:"not null" json:"-"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ErrResponce struct {
	ErrMessage string
}

type ResponceData struct {
	Message string
	Data    interface{}
}
