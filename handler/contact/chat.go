package contact

import (
	"encoding/json"
	"net/http"

	"go-chatapp/service"
	ws "go-chatapp/websocket"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	manager *ws.Manager
	service *service.ChatService
}

func NewWebSocketHandler(
	manager *ws.Manager,
	service *service.ChatService,
) *WebSocketHandler {
	return &WebSocketHandler{
		manager: manager,
		service: service,
	}
}

type ChatMessage struct {
	ConversationID int64  `json:"conversation_id"`
	Content        string `json:"content"`
}

var upgrader = gorilla.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebSocketHandler) Connect(c *gin.Context) {

	userID := c.GetInt("user_id")

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	h.manager.Add(userID, conn)

	defer func() {
		h.manager.Remove(userID)
		conn.Close()
	}()

	for {

		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var input ChatMessage

		err = json.Unmarshal(data, &input)
		if err != nil {
			conn.WriteJSON(gin.H{
				"error": "invalid message format",
			})
			continue
		}

		message, err := h.service.SendMessage(
			c.Request.Context(),
			int64(userID),
			input.ConversationID,
			input.Content,
		)

		if err != nil {
			conn.WriteJSON(gin.H{
				"error": err.Error(),
			})
			continue
		}

		conn.WriteJSON(message)
	}
}
