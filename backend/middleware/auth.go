package middleware

import "github.com/gin-gonic/gin"

// Fake rol de admin
func FakeAdmin() gin.HandlerFunc {
	// Inyecto en el contexto un user admin
	return func(c *gin.Context) {
		// Simulo un userID y rol admin
		c.Set("userID", "demo-user")
		c.Set("userRole", "admin")
		c.Next()
	}
}
