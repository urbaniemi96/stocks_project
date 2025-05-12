package main

import (
  "log"
  "github.com/gin-gonic/gin"
  //"fmt"
)


//Se ejecuta al hacer el go run
func main() {
  // Conecto a la DB
  initDB()
  // Inicializo gin
  r := gin.Default()

  r.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
  
    if c.Request.Method == "OPTIONS" {
      c.AbortWithStatus(204)
      return
    }
  
    c.Next()
  })

  // Rutas hacia los manejadores
  //r.GET("/fetch", getStocksHandler) // Traigo datos de la API y guardo en la DB
  r.GET("/fetch", StartFetchHandler) // Traigo datos de la API y guardo en la DB
  r.GET("/stocks", listStocksHandler) // Muestro los datos guardados en la DB
  r.GET("/enrich", StartEnrichHandler) // Muestro los datos guardados en la DB
  r.GET("/task/:id", FetchStatusHandler) //Consulto el estado de la tarea con :id
  
  //r.GET("/recommend", recommendHandler) // Recomiendo mejor stock

  // Arranco el servidor en el puerto 8080
  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }
}