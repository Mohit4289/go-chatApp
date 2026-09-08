package acc

import (
	"go-chatapp/service"

	"github.com/gin-gonic/gin"
)

type LoginUser struct {
	Email    string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginHandler struct {
	userService *service.UserService
}

func LoginUserHandler(userService *service.UserService) *LoginHandler {
	return &LoginHandler{userService: userService}
}

func (s *LoginHandler) Login(ctx *gin.Context) {

}
