package utils

import (
	"fmt"
	"net/smtp"

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

func SendEmailWarning(email, newIP string) (string, string) {
	// адрес электронной почты отправителя
	senderMail := ""
	passwordSenderMail := ""

	//адрес электронной почты получателя
	userEmail := email
	to := []string{userEmail}

	// данные для настройки SMTP-клиента
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	// содержание письма
	subject := "Warning, logging in from a new IP address\n\n"
	body := fmt.Sprintf("IP address: %s", newIP)
	message := []byte(subject + body)

	// для аутентификации
	auth := smtp.PlainAuth("", senderMail, passwordSenderMail, host)

	// отправка через SMTP-клиент
	err := smtp.SendMail(address, auth, senderMail, to, message)
	if err != nil {
		fmt.Println(err)
		return "", err.Error()
	}

	return "email warning was successfully sent", ""
}
