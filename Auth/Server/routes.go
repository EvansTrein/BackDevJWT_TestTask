package server

import (
	"AuthServ/handlers"

	"github.com/gin-gonic/gin"
)

func InitRotes() {
	router := gin.Default()

	router.POST("/auth/:guid", handlers.AuthHandler)

	router.POST("/auth/refresh", handlers.AuthRefreshHandler)

	router.DELETE("/auth/delSession/:guid", handlers.DelSessionHandler)


	router.POST("/userCreate", handlers.UserCreateHandler)

	router.GET("/user/:guid", handlers.Middleware, handlers.UserHandler)

	router.PUT("/user/:guid/update", handlers.UserUpdateHandler)

	router.DELETE("/user/:guid/del", handlers.UserDelHandler)


	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":4000")
}
