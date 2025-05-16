package main

import (
	//"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

var db *gorm.DB

// Inicio la conexión a la DB, creo las tablas si no existen
func initDB() {
	var err error
	dsn := getDBDSN() 
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la db: %v", err)
	}
	// Creo tablas (si no existen) de stocks, historial, y recomendaciones
	db.AutoMigrate(&Stock{}, &HistoricalPoint{}, &Recommendation{})
	if err != nil {
		log.Fatalf("Error al crear tablas: %v", err)
	}
}

func saveStocks(stocks []Stock) error {
	// Inicio transacción con rollback automático en caso de error
	return db.Transaction(func(tx *gorm.DB) error {
		// Recorro los stocks
		for _, s := range stocks {
			// Si ya existe el ticker , actualiza los campos
			if err := tx.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&s).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

