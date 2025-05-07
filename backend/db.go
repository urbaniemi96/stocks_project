package main

import (
  "context"
  "log"
  "github.com/jackc/pgx/v5"
)

var db *pgx.Conn

// Inicio la conexión a la DB, creo la tabla stocks si no existe
func initDB() {
  var err error
  url := getDBURL() // En config.go
  db, err = pgx.Connect(context.Background(), url) // Paso un contexto raíz vacío (para pruebas) porque no tengo un contexto superior (relacionado a timeouts y cancelaciones, investigar más)
  if err != nil {
    log.Fatalf("Error al conectar a la db: %v", err)
  }
  // Creo tabla si no existe
  _, err = db.Exec(context.Background(), `
    CREATE TABLE IF NOT EXISTS stocks (
      ticker TEXT PRIMARY KEY,
      company TEXT,
      target_from FLOAT,
      target_to FLOAT,
      action TEXT,
      brokerage TEXT,
      rating_from TEXT,
      rating_to TEXT,
      time TIMESTAMPTZ
    );
  `)
  if err != nil {
    log.Fatalf("Error al crear la tabla: %v", err)
  }
}

func saveStocks(ctx context.Context, stocks []Stock) error {
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
  // Al hacer el commit, el rollback del defer queda sin efecto (el commit primero se hace, y DESPUES se retorna su valor (por más que que defer se ejecute antes de retornar))
  return tx.Commit(ctx)
}