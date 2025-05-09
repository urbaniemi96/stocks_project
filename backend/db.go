package main

import (
  //"context"
  "log"
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
  "gorm.io/gorm/clause"
)

var db *gorm.DB

// Inicio la conexión a la DB, creo la tabla stocks si no existe
func initDB() {
  var err error
  dsn := getDBDSN() // En config.go
  db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    log.Fatalf("Error al conectar a la db: %v", err)
  }
  // Creo tabla si no existe usando AutoMigrate
  db.AutoMigrate(&Stock{})
  if err != nil {
    log.Fatalf("Error al crear la tabla: %v", err)
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

/*func saveStocks(ctx context.Context, stocks []Stock) error {
  // Inicio transacción (con contexto raiz vacío)
  tx, err := db.Begin(ctx)
  if err != nil {
	  return err
  }
  // Preparo un rollback en caso de que la transacción falle
  defer tx.Rollback(ctx)

  for _, s := range stocks {
	// Ejecuto las consultas con placeholder para los parámetros
    _, err := tx.Exec(ctx,
	`UPSERT INTO stocks (
		ticker, company, target_from, target_to, action,
		brokerage, rating_from, rating_to, time
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
	s.Ticker, s.Company, s.TargetFrom, s.TargetTo,
	s.Action, s.Brokerage, s.RatingFrom, s.RatingTo, s.Time,
	)
    if err != nil {
	    return err
	  }
  }
  retur := tx.Commit(ctx)

  // Al hacer el commit, el rollback del defer queda sin efecto (el commit primero se hace, y DESPUES se retorna su valor (por más que que defer se ejecute antes de retornar))
  return retur
}*/