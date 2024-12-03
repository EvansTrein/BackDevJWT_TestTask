package server

import (
	"AuthServ/handlers"

	_ "AuthServ/docs"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

func InitRotes() {
	router := gin.Default()

	router.POST("/auth/:guid", handlers.AuthHandler)

	router.POST("/auth/refresh", handlers.AuthRefreshHandler)

	router.DELETE("/auth/delSession/:guid", handlers.DelSessionHandler)


	router.POST("/userCreate", handlers.UserCreateHandler)

	router.GET("/user/:guid", handlers.Middleware, handlers.UserHandler)

	router.PUT("/user/:guid/update", handlers.Middleware, handlers.UserUpdateHandler)

	router.DELETE("/user/:guid/del", handlers.Middleware, handlers.UserDelHandler)
	

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":4000")
}
