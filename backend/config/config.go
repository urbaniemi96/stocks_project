package config

import (
	"github.com/joho/godotenv"
	"runtime"
	"log"
	"os"
	"path/filepath"
)

// Todas las funciones init() dentro del paquete main se ejecutan ANTES del main(). Ideal para cargar configuraciones
func init() {
	// Obtengo la ruta de este archivo (para no depender de dónde se ejecuta el go run)
    _, thisFile, _, ok := runtime.Caller(0)
    if !ok {
        log.Fatal("No se pudo obtener la ruta de config.go")
    }
		// Subo a /backend
    projectBackend := filepath.Join(filepath.Dir(thisFile), "..")
    envPath := filepath.Join(projectBackend, ".env")
	// Cargo .env
	err := godotenv.Load(envPath)
	// Detecto error y detenco ejecución
	if err != nil {
		log.Fatal("Error al cargar el archivo .env - ERROR: ", err)
	}
}

// Obtengo la dsn de la db desde el .env (para GORM)
func GetDBDSN() string {
	return os.Getenv("DATABASE_DSN")
}

// Obtengo URL de la API
func GetAPIURL() string {
	return os.Getenv("API_URL")
}

// Obtengo KEY de la API
func GetAPIKEY() string {
	return os.Getenv("API_KEY")
}
