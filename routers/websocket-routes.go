package routers

import (
	"go-chatapp/handler/chat"
	"go-chatapp/middleware"

	"github.com/gin-gonic/gin"
)

func SetupWebSocketRoutes(r *gin.Engine, websocketHandler *chat.WebSocketHandler) {
	chatGroup := r.Group("/api/chat")
	{
		chatGroup.GET("/chat", middleware.TokenVerficationMiddleware(), websocketHandler.Connect)
	}
}
