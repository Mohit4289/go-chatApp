package contact

import (
	"go-chatapp/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContactListHandler struct {
	contactService *service.ContactService
}

func ContactUserListHandler(contactService *service.ContactService) *ContactListHandler {
	return &ContactListHandler{contactService: contactService}
}

func (h *ContactListHandler) UserList(ctx *gin.Context) {
	userData, err := h.contactService.GetAllUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "not able to fetch contact list",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "fetched contact list",
		"data":    userData,
	})
}
