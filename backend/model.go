package main

import "time"

// Creo el modelo de acciones
type Stock struct {
  Ticker     string    `json:"ticker"`
  Company    string    `json:"company"`
  TargetFrom float64   `json:"target_from"`
  TargetTo   float64   `json:"target_to"`
  Action     string    `json:"action"`
  Brokerage  string    `json:"brokerage"`
  RatingFrom string    `json:"rating_from"`
  RatingTo   string    `json:"rating_to"`
  Time       time.Time `json:"time"`
}
