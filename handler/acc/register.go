package acc

import (
	"net/http"

	"go-chatapp/service"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var userdata User
	err := ctx.ShouldBindBodyWithJSON(&userdata)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "please send data",
		})
		return
	}

	createdUser, err := h.userService.RegisterAcc(ctx, service.User{
		Name:     userdata.Name,
		Email:    userdata.Email,
		Password: userdata.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
		"user":    createdUser,
	})

}
