package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const JWT_SECRET = "adsiAWqegd234123Sgke"

func GenerateAcessToken(guid, clientIP string) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		// создаем Payload токена
		"guid":     guid,
		"clientIP": clientIP,
		"exp":      time.Now().Add(time.Minute * 1).Unix(), // время жизни токена
	})

	// подписываем токен
	signedAccessToken, err := accessToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	return token, nil

}

func GenerateRefreshToken() (string, error) {
	data := make([]byte, 32)
	// генерируем случайные данные для токена
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	// приводим токен к формату base64
	refreshToken := base64.StdEncoding.EncodeToString(data)

	return refreshToken, nil
}
