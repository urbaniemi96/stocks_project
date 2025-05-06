package main

import (
  "context"
  "log"
  "github.com/jackc/pgx/v5"
)


// Inicio la conexión a la DB, creo la tabla stocks si no existe
func initDB() {
  var db *pgx.Conn
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