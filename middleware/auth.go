// middleware/auth.go
package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
