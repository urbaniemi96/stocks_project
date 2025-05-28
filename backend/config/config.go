package config

import (
	"github.com/joho/godotenv"
	"runtime"
	"log"
	"os"
	"path/filepath"
	"strings"

	"sync"
)

// Para cargar solo una vez el .env
var loadEnvOnce sync.Once

// Todas las funciones init() dentro del paquete main se ejecutan ANTES del main(). Ideal para cargar configuraciones -- CAMBIADO POR RECOMENDACIÓN DEL RECLUTADOR
func InitEnv() {
	// Detecto si esto es un test y salto la carga del .env
	if runningTests() {
			return 
	}
	// Cargo solo una vez el .env
		loadEnvOnce.Do(func() {
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
			// Detecto error y detengo ejecución
			if err != nil {
				log.Fatal("Error al cargar el archivo .env - ERROR: ", err)
			}
		})
}

// Detecto en los argumentos si se está ejecutando un test
func runningTests() bool {
    for _, arg := range os.Args {
        if strings.HasPrefix(arg, "-test.") {
            return true
        }
    }
    return false
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
