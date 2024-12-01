package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Middleware(ctx *gin.Context) {

	log.Println("Middleware")
	
	log.Println(ctx.Request.Header)

	// Создаем новый HTTP-запрос
	req, err := http.NewRequestWithContext(ctx.Request.Context(), "POST", "http://localhost:4000/auth/refresh", ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Добавляем необходимые заголовки
	req.Header.Set("Authorization", "Bearer your_token_here")
	req.Header.Set("RefreshToken", "custom_value")

	// Отправляем запрос на внешний сервис
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	log.Println(resp)
	log.Println(resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response body"})
		return
	}

	log.Println(data)
	log.Println(data["newTokens"])


	log.Println("Middleware after")

	ctx.Next()
	// ctx.Abort()
}
