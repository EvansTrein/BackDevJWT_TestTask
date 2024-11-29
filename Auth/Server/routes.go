package server

import (
	"github.com/gin-gonic/gin"
)

func InitRotes() {
	router := gin.Default()

	router.POST("/auth/:guid", AuthHandler)

	router.POST("/auth/refresh", AuthRefreshHandler)

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":4000")
}
