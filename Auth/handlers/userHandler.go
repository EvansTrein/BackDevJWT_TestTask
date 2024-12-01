package handlers

import (
	"github.com/gin-gonic/gin"
)

func UserHandler(ctx *gin.Context) {

	ctx.JSON(200, gin.H{"message": "data"})
}
