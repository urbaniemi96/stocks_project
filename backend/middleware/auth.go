package middleware

import "github.com/gin-gonic/gin"

func FakeAdmin() gin.HandlerFunc {
  return func(c *gin.Context) {
    // simulamos un userID y rol
    c.Set("userID", "demo-user")
    c.Set("userRole", "admin")
    c.Next()
  }
}