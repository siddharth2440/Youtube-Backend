package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/youtube/utils"
)

func ChkAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("authCookie")
		if err != nil || cookie == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		// VerifyToken
		token, err := utils.Verify_JWT_Token(cookie)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok && !token.Valid {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		ctx.Set("userId", claims["userId"])
		ctx.Next()
	}
}
