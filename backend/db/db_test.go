package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/urbaniemi96/stocks_project/backend/model"
)

func TestSaveStocks(t *testing.T) {
	// Creo mock sql
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	// Abro GORM sobre CockroachDB (Postgres) con la conexión mockeada
	dialector := postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	// Reasigno variable global para que use el mock
	DB = gormDB

	// Modelo a probar
	stocks := []model.Stock{
		{
			Ticker:     "AAPL",
			Company:    "Apple Inc.",
			TargetFrom: 150.0,
			TargetTo:   160.0,
			Action:     "buy",
			Brokerage:  "BrokerX",
			RatingFrom: "B",
			RatingTo:   "A",
		},
	}

	// Indico que espere inicio de transacción
	mock.ExpectBegin()
	// Espero que la consulta inicie con un insert (puedo testear otras partes de la consulta aquí pero me está dando problemas)
	mock.ExpectExec(`INSERT INTO "stocks"`).
		WithArgs(
			stocks[0].Ticker,
			stocks[0].Company,
			stocks[0].TargetFrom,
			stocks[0].TargetTo,
			stocks[0].Action,
			stocks[0].Brokerage,
			stocks[0].RatingFrom,
			stocks[0].RatingTo,
			sqlmock.AnyArg(), // por ejemplo CreatedAt/Time
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	// Espero commit
	mock.ExpectCommit()

	/*
		// Puedo esperar un rollback forzando error de la siguiente forma
		mock.ExpectExec(`INSERT INTO "stocks"`).
			WillReturnError(fmt.Errorf("violación de unique"))
		mock.ExpectRollback()
	*/

	// Llamo la función
	err = SaveStocks(stocks)
	require.NoError(t, err)

	// Verifico que todas las expectativas se hayan cumplido
	require.NoError(t, mock.ExpectationsWereMet())
}
