package middleware

import "github.com/gin-gonic/gin"

// Fake rol de admin
func FakeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulo un userID y rol admin
		c.Set("userID", "demo-user")
		c.Set("userRole", "admin")
		c.Next()
	}
}
