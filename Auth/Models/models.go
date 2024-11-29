package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	GUID     string `gorm:"not null"`
	NameUser string `gorm:"not null"`
	EmaiUser string `gorm:"not null;unique"`
}

type ClientSession struct {
	gorm.Model
	RefreshToken string `json:"-"`
	AdressIp     string 
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
}
