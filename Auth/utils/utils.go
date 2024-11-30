package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Hashing(s string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckHashing(s, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	return err == nil
}
