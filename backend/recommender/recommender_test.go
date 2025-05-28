package recommender

import (
	"regexp"
	"math"
	"testing"
	"time"
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/urbaniemi96/stocks_project/backend/db"

	"github.com/urbaniemi96/stocks_project/backend/tasks"

)
// (Comentario de la IA al preguntar por require.InEpsilon. Me pareció informativo dejarlo)
// floatMatcher implementa sqlmock.Argument para comparar floats con tolerancia
// no usamos require.InEpsilon aquí porque necesitamos una interfaz Argument
// que sqlmock pueda invocar al hacer match de los parámetros en la consulta.
// require.InEpsilon funciona sólo dentro de un Test para aserciones posteriores,
// pero no sirve para interceptar y validar los valores pasados a Exec.
// - expected: valor esperado
// - tol: tolerancia máxima permitida
type floatMatcher struct{
	expected float64
	tol      float64
}
// Verifico si el valor recibido "v" es float64 y si difiere del esperado con la tolerancia
func (f floatMatcher) Match(v driver.Value) bool {
	// Intento convertirlo a float64
	val, ok := v.(float64)
	if !ok { 
		return false 
	}
	// Calculo la diferencia entre valores y si entra dentro de la tolerancia
	return math.Abs(val - f.expected) <= f.tol
}
// Testeo RecalculateRecommendations usando sqlmock
func TestRecalculateRecommendations_SQLMock(t *testing.T) {
	// Configuración inicial
	fakeTaskID := "test-task-123"
	tasks.TasksMu.Lock()
	tasks.Tasks[fakeTaskID] = &tasks.TaskInfo{Status: "in-progress", PagesFetched: 0, Error: ""}
	tasks.TasksMu.Unlock()

	// Inicializo sqlmock y GORM
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	// Abro GORM sobre CockroachDB (Postgres) con la conexión mockeada
	dialector := postgres.New(postgres.Config{
		Conn: sqlDB, 
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	// Reasigno variable global para que use el mock
	db.DB = gormDB

	// Mockeo lectura de puntos históricos (defino las columnas que devolverá y añado dos filas de ejemplo)
	histRows := sqlmock.NewRows([]string{"ticker", "date", "open", "high", "low", "close"}).
		AddRow("AAA", time.Now(), 100.0, 115.0, 95.0, 110.0).
		AddRow("AAA", time.Now(), 200.0, 220.0, 180.0, 210.0)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"historical_points\"")).
		WillReturnRows(histRows)

	// Mockeo consulta de rating con 2 args (ticker + limit(cualquier valor))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT \"rating_to\" FROM \"stocks\" WHERE ticker = $1 ORDER BY \"stocks\".\"ticker\" LIMIT $2")).
		WithArgs("AAA", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"rating_to"}).AddRow("buy"))

	// Mockeo inserción de recomendación con ON CONFLICT 
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO \"recommendations\"")).
		WithArgs(
			"AAA",
			floatMatcher{expected: (0.075/(0.2+1e-6))*1.2, tol: 1e-6},
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Ejecuto función a testear
	err = RecalculateRecommendations(fakeTaskID)
	// Verifico que se hayan incrementado las páginas procesadas
	tasks.TasksMu.Lock()
	if tasks.Tasks[fakeTaskID].PagesFetched != 1 {
			t.Errorf("expected PagesFetched to be 1, got %d", tasks.Tasks[fakeTaskID].PagesFetched)
	}
	tasks.TasksMu.Unlock()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
