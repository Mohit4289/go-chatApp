package contact

import (
	"go-chatapp/service"

	"github.com/gin-gonic/gin"
)

type ContactListHandler struct {
	contactService *service.UserService
}

func ContactUserListHandler(userService *service.UserService) *ContactListHandler {
	return &ContactListHandler{contactService: userService}
}

func (h *ContactListHandler) UserList(ctx *gin.Context) {

}
