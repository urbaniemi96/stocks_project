package main

import "time"

// Creo el modelo de acciones
type Stock struct {
	Ticker     string    `json:"ticker" gorm:"primaryKey"` // Clave primaria
	Company    string    `json:"company"`
	TargetFrom float64   `json:"target_from"`
	TargetTo   float64   `json:"target_to"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
}

// Creo modelo de histórico de acciones con los precios en distintas fechas (para enriquecer los datos)
type HistoricalPoint struct {
	Ticker string    `gorm:"primaryKey"`
	Date   time.Time `gorm:"primaryKey"`
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

// Creo el modelo de recomendaciones para las acciones
type Recommendation struct {
    Ticker    string    `gorm:"primaryKey"`
    Score     float64   `gorm:"not null"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}