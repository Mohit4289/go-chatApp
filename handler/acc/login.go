package acc

import (
	"errors"
	"go-chatapp/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginHandler struct {
	userService *service.UserService
}

func LoginUserHandler(userService *service.UserService) *LoginHandler {
	return &LoginHandler{userService: userService}
}

func (s *LoginHandler) Login(ctx *gin.Context) {
	var userData LoginUser
	err := ctx.ShouldBindBodyWithJSON(&userData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "dont send blank body",
		})
		return
	}

	checkingData, accessToken, refreshToken, err := s.userService.LoginAcc(ctx, service.LoginData{
		Email:    userData.Email,
		Password: userData.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrUserNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "failed to verify user",
				"err":     err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to verify user",
			"err":     err.Error(),
		})
		return
	}

	ctx.SetCookie(
		"access_token",
		accessToken,
		60*60*24,
		"/",
		"",
		false,
		true,
	)
	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		60*60*24,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Logined successfull",
		"verfied": checkingData,
	})
}
