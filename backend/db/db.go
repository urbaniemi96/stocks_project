package db

import (
	//"context"
	"github.com/urbaniemi96/stocks_project/backend/config"
	"github.com/urbaniemi96/stocks_project/backend/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

var DB *gorm.DB

// Inicio la conexión a la DB, creo las tablas si no existen
func InitDB() {
	var err error
	// Cargo el archivo .env si no está cargado ya
	config.InitEnv()
	dsn := config.GetDBDSN()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la db: %v", err)
	}
	// Creo tablas (si no existen) de stocks, historial, y recomendaciones
	DB.AutoMigrate(&model.Stock{}, &model.HistoricalPoint{}, &model.Recommendation{})
	if err != nil {
		log.Fatalf("Error al crear tablas: %v", err)
	}
}

func SaveStocks(stocks []model.Stock) error {
	// Inicio transacción con rollback automático en caso de error
	return DB.Transaction(func(tx *gorm.DB) error {
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
