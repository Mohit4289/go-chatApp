package routers

import (
	"go-chatapp/handler/acc"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, userHandler *acc.UserHandler) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", userHandler.Register)
	}
}
