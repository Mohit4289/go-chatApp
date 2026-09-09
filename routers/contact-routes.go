package routers

import (
	"go-chatapp/handler/contact"

	"github.com/gin-gonic/gin"
)

func SetupContactRoutes(r *gin.Engine, listUserHandler *contact.ContactListHandler) {
	contactGroup := r.Group("/api/contact")
	{
		contactGroup.POST("/contact-user", listUserHandler.UserList)
	}
}
