package main

import (
    "github.com/gin-gonic/gin"
)


//Se ejecuta al hacer el go run
func main() {

    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.Run() // default :8080
}