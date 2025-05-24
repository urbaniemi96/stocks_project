package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Configuro router en modo test y aplico middleware que le paso
func makeTestRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// Modo test
	gin.SetMode(gin.TestMode)
	// Creo contexto gin
	r := gin.New()
	// Uso middlewares pasados
	r.Use(middlewares...)
	return r
}
// Testeo el seteo de userID y de userRole
func TestFakeAdmin_SetsUserIDAndRole(t *testing.T) {
	// Creo contexto con el middleware para crear el admin
	r := makeTestRouter(FakeAdmin())
	// Ruta dummy que lee del Context
	r.GET("/check", func(c *gin.Context) {
		uid, exists := c.Get("userID")
		// Testeo que exista el userID y sea demo-user
		require.True(t, exists, "debe existir userID en el Context")
		require.Equal(t, "demo-user", uid)
		// Testeo que exista el userRole y sea admin
		role, exists := c.Get("userRole")
		require.True(t, exists, "debe existir userRole en el Context")
		require.Equal(t, "admin", role)

		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	// Solicitud a la ruta dummy
	req := httptest.NewRequest("GET", "/check", nil)
	r.ServeHTTP(w, req)
	// Checkeo status OK 200
	require.Equal(t, http.StatusOK, w.Code)
}
// Testeo el RequireAdmin cuando el usuario es un admin
func TestRequireAdmin_AllowsWhenAdmin(t *testing.T) {
	// Inyecto en el contexto el user admin antes de RequireAdmin
	r := makeTestRouter(
		FakeAdmin(),
		RequireAdmin(),
	)
	// Ruta solo accesible por el admin
	r.GET("/secret", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	// Solicitud a la ruta de admin y aplico los middlewares
	r.ServeHTTP(w, httptest.NewRequest("GET", "/secret", nil))
	// Checkeo status OK 200
	require.Equal(t, http.StatusOK, w.Code)
}
// Testeo el RequireAdmin cuando el usuario NO es un admin
func TestRequireAdmin_ForbiddenWhenNotAdmin(t *testing.T) {
	// Inyecto en el contexto userRole distinto de "admin"
	r := makeTestRouter(
		func(c *gin.Context) {
			c.Set("userRole", "user")
			c.Next()
		},
		RequireAdmin(),
	)
	// Ruta solo accesible por el admin
	r.GET("/secret", func(c *gin.Context) {
		t.Fatal("handler no debería ejecutarse cuando userRole != admin")
	})

	w := httptest.NewRecorder()
	// Solicitud a la ruta de admin y aplico los middlewares
	r.ServeHTTP(w, httptest.NewRequest("GET", "/secret", nil))
	// Checkeo status forbidden 403
	require.Equal(t, http.StatusForbidden, w.Code)
	require.JSONEq(t, `{"error":"forbidden"}`, w.Body.String())
}
// Testeo el RequireAdmin cuando no hay usuario
func TestRequireAdmin_ForbiddenWhenNoRoleSet(t *testing.T) {
	// No se inyecta ningún rol
	r := makeTestRouter(RequireAdmin())
	// Ruta solo accesible por el admin
	r.GET("/secret", func(c *gin.Context) {
		t.Fatal("handler no debería ejecutarse cuando no hay userRole")
	})

	w := httptest.NewRecorder()
	// Solicitud a la ruta de admin y aplico los middlewares
	r.ServeHTTP(w, httptest.NewRequest("GET", "/secret", nil))
	// Checkeo status forbidden 403
	require.Equal(t, http.StatusForbidden, w.Code)
	require.JSONEq(t, `{"error":"forbidden"}`, w.Body.String())
}
