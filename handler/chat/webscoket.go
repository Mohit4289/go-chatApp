package chat

import (
	"net/http"

	"go-chatapp/websocket"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

var upgrader = gorilla.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	manager *websocket.Manager
}

func NewWebSocketHandler(
	manager *websocket.Manager,
) *WebSocketHandler {

	return &WebSocketHandler{
		manager: manager,
	}
}

func (h *WebSocketHandler) Connect(ctx *gin.Context) {
	var userID = ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "user not found",
		})
		return
	}

	conn, err := upgrader.Upgrade(
		ctx.Writer,
		ctx.Request,
		nil,
	)
	if err != nil {
		return
	}

	h.manager.Add(userID, conn)

	defer func() {
		h.manager.Remove(userID)
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()

		if err != nil {
			break
		}
	}
}
