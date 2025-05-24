package detail

import (
    "net/http"
    "net/url"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/require"
    "github.com/urbaniemi96/stocks_project/backend/model"
)

// Simulo un contexto 
func makeContext(rawQuery string) *gin.Context {
    req := &http.Request{URL: &url.URL{RawQuery: rawQuery}}
    return &gin.Context{Request: req}
}

// Test de filtros por defecto
func TestParseHistoryFilters_Defaults(t *testing.T) {
    c := makeContext("")
    hf, err := ParseHistoryFilters(c)
    require.NoError(t, err)
    // Valores por defecto
    require.Equal(t, 90, hf.Days)
		require.Nil(t, hf.StartDate, "StartDate debería ser nil")
    require.Nil(t, hf.EndDate,   "EndDate debería ser nil")
    require.Nil(t, hf.MinPrice,  "MinPrice debería ser nil")
    require.Nil(t, hf.MaxPrice,  "MaxPrice debería ser nil")
    require.Nil(t, hf.MinVolume, "MinVolume debería ser nil")
    require.False(t, hf.OrderDesc, "OrderDesc debe ser falso por defecto (asc)")
}

// Test de filtros con parámetros
func TestParseHistoryFilters_WithParams(t *testing.T) {
	// Valores a testear
    vals := url.Values{
        "days":       {"10"},
        "start_date": {"2025-01-01"},
        "end_date":   {"2025-01-31"},
        "min_price":  {"1.5"},
        "max_price":  {"2.5"},
        "min_volume": {"1000"},
        "order":      {"desc"},
    }
		// Creo el contexto con los valores a testear
    c := makeContext(vals.Encode())
    hf, err := ParseHistoryFilters(c)
    require.NoError(t, err)

    require.Equal(t, 10, hf.Days)
		// Creo objetos de las fechas a testear
    sd := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
    ed := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)
		// Comparo el filtro con la fecha a testear
    require.True(t, hf.StartDate.Equal(sd), "StartDate erronea")
    require.True(t, hf.EndDate.Equal(ed), "EndDate erronea")
		// Comparo el resto de valores (checkeando punteros no nulos por cómo uso el código en ParseHistoryFilters)
    require.NotNil(t, hf.MinPrice)
    require.Equal(t, 1.5, *hf.MinPrice, "MinPrice erroneo")

    require.NotNil(t, hf.MaxPrice)
    require.Equal(t, 2.5, *hf.MaxPrice, "MaxPrice erroneo")

    require.NotNil(t, hf.MinVolume)
    require.Equal(t, int64(1000), *hf.MinVolume, "MinVolume erroneo")

    require.True(t, hf.OrderDesc, "OrderDesc debería ser true cuando order=desc")
}

// Test del cálculo de riesgo y recompensa
func TestCalcRiskReward(t *testing.T) {
	// Creo modelo de históricos a probar
    hist := []model.HistoricalPoint{
        {
            Date:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
            Open:  10, High: 12, Low: 9, Close: 11,
        },
        {
            Date:  time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
            Open:  20, High: 22, Low: 18, Close: 19,
        },
    }

    rr := CalcRiskReward(hist)
    // Testeo labels
    require.Len(t, rr.Labels, 2) // longitud
    require.Equal(t, "2025-01-01", rr.Labels[0]) // valor
    require.Equal(t, "2025-01-02", rr.Labels[1]) // valor

		// Testeo con cierta tolerancia en los valores (epsilon)
    // testeo volatilidades: (High-Low)/Open*100
    require.InEpsilon(t, 30.0, rr.Volatilities[0], 1e-6)
    require.InEpsilon(t, 20.0, rr.Volatilities[1], 1e-6)

    // Testeo potenciales: (Close-Open)/Open*100
    require.InEpsilon(t,  10.0, rr.Potentials[0], 1e-6)
    require.InEpsilon(t,  -5.0, rr.Potentials[1], 1e-6)
}