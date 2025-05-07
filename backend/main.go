package main

import (
  "log"
  "github.com/gin-gonic/gin"
)


//Se ejecuta al hacer el go run
func main() {
  // Conecto a la DB
  initDB()
  // Inicializo gin
  r := gin.Default()
  // Rutas hacia los manejadores
  r.GET("/fetch", getStocksHandler) // Traigo datos de la API y guardo en la DB
  r.GET("/stocks", listStocksHandler) // Muestro los datos guardados en la DB
  r.GET("/recommend", recommendHandler) // Recomiendo mejor stock

  // Arranco el servidor en el puerto 8080
  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }
}