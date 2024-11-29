package server

import (
	"AuthServ/handlers"

	"github.com/gin-gonic/gin"
)

func InitRotes() {
	router := gin.Default()

	router.POST("/auth/:guid", handlers.AuthHandler)

	router.POST("/auth/refresh", handlers.AuthRefreshHandler)

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":4000")
}
