package main

import (
	"github.com/gin-gonic/gin"
	"github.com/urbaniemi96/stocks_proyect/backend/middleware"
	"log"
	//"fmt"
)

// Se ejecuta al hacer el go run
func main() {
	// Conecto a la DB
	initDB()
	// Inicializo gin
	r := gin.Default()

	// Inyecto "usuario admin"
	r.Use(middleware.FakeAdmin())

	// Solo para desarrollo permito todo
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
	
	// Muestro los stocks guardados en la DB
	r.GET("/stocks", listStocksHandler)    
	
	// Traigo los históricos de un ticker
	r.GET("/stocks/:ticker/detail", StockDetailHandler)

	// Muestro el Top 20 recomendaciones
	r.GET("/recommendations/top20", TopRecommendationsHandler)

	// Rutas del rol admin”
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAdmin())
	admin.POST("/recalculate", RecalculateRecommendationsHandler) // Recalculo score de recomendación con los datos guardados
	admin.GET("/fetch", StartFetchHandler)     // Traigo datos de la API del desafío y guardo en la DB
	admin.GET("/enrich", StartEnrichHandler)   // Inicio enriquecimiento de datos con Yahoo Finance
	admin.GET("/task/:id", FetchStatusHandler) //Consulto el estado de la tarea con :id


	// Arranco el servidor en el puerto 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
