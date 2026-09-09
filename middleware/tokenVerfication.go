package middleware

import (
	"go-chatapp/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TokenVerficationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "failed to get cookie",
			})
			ctx.Abort()
			return
		}

		secret := config.EnvConfig().JWT_SECRET
		if secret == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error: JWT secret not configured",
			})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return []byte(secret), nil
		})
		if err != nil && !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			ctx.Abort()
			return
		}

		userIDVal, ok := claims["user_id"]
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "user id not exsist",
			})
			ctx.Abort()
			return
		}

		var id int
		switch v := userIDVal.(type) {
		case float64:
			id = int(v)
		case int:
			id = v
		case int64:
			id = int(v)
		default:
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user_id format in token",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", id)
		ctx.Next()
	}

}
