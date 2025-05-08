package main

import (
  "context"
  "net/http"
  "github.com/gin-gonic/gin"
  //"fmt"
)
// Manejador para obtener los stocks
func getStocksHandler(c *gin.Context) {
  // Traigo los stocks
  stocks, err := fetchAllStocks()

  if err != nil {
    c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
    return
  }

  // Guardo los stocks en la DB
  if err := saveStocks(context.Background(), stocks); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  // Devuelvo todo Ok
  c.JSON(http.StatusOK, stocks)
}

// Manejador para listar los stocks guardados
func listStocksHandler(c *gin.Context) {
  rows, _ := db.Query(context.Background(), `SELECT ticker, company, target_from, target_to, action, brokerage, rating_from, rating_to, time FROM stocks`)
  defer rows.Close()

  var all []Stock
  for rows.Next() {
    var s Stock
    rows.Scan(&s.Ticker, &s.Company, &s.TargetFrom, &s.TargetTo,
      &s.Action, &s.Brokerage, &s.RatingFrom, &s.RatingTo, &s.Time)
    all = append(all, s)
  }
  c.JSON(http.StatusOK, all)
}

// Manejador para recomendar mejores stocks (REVISAR JUNTO CON EL ALGORITMO)
func recommendHandler(c *gin.Context) {
  best, err := recommendBest(context.Background())
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, best)
}