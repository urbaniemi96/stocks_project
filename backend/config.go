package main

import (
    "github.com/joho/godotenv"
    "log"
    "os"
    //"fmt"
)

// Todas las funciones init() dentro del paquete main se ejecutan ANTES del main(). Ideal para cargar configuraciones
func init() {
	// Cargo .env
    err := godotenv.Load()
	// Detecto error y detenco ejecución
    if err != nil {
        log.Fatal("Error al cargar el archivo .env - ERROR: ", err)
    }
}

// Obtengo la url de la db desde el .env
func getDBURL() string {
	return os.Getenv("DATABASE_URL")
}
// Obtengo URL de la API
func getAPIURL() string {
	return os.Getenv("API_URL")
}
// Obtengo KEY de la API
func getAPIKEY() string {
	return os.Getenv("API_KEY")
}