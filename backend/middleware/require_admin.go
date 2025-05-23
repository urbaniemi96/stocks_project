package middleware

import (
  "net/http"
  "github.com/gin-gonic/gin"
)
// Middleware para requerir rol admin
func RequireAdmin() gin.HandlerFunc {
  return func(c *gin.Context) {
    role, _ := c.Get("userRole")
    if role != "admin" {
      c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
      return
    }
    c.Next()
  }
}