package middleware

import (
    "context"
    "net/http"

    "homefinder/db"

    "github.com/gin-gonic/gin"
)

func OwnerRoleMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.MustGet("user_id").(string)

        var role string
        err := db.Conn.QueryRow(context.Background(), "SELECT role FROM users WHERE user_id=$1", userID).Scan(&role)
        if err != nil || role != "homeowner" {
            c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
            c.Abort()
            return
        }

        c.Next()
    }
}
