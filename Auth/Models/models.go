package models

import "gorm.io/gorm"

type ClientSession struct {
	gorm.Model
	GUID         string `gorm:"not null"`
	NameUser     string `gorm:"not null"`
	RefreshToken string
	Emai         string
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
