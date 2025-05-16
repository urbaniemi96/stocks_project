package main

import (
	//"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
	"log"
)

type HistoryFilters struct {
    Days      int
    StartDate *time.Time
    EndDate   *time.Time
    MinPrice  *float64
    MaxPrice  *float64
    MinVolume *int64
    OrderDesc bool
}

type StockDetailResponse struct {
    Stock              Stock                   `json:"stock"`
    History            []HistoricalPoint       `json:"history"`
    RiskReward         RiskRewardData          `json:"riskReward"`
    RatingDistribution map[string]int          `json:"ratingDistribution"`
}

type RiskRewardData struct {
    Labels       []string  `json:"labels"`
    Volatilities []float64 `json:"volatilities"`
    Potentials   []float64 `json:"potentials"`
}

// Inicia una goroutine que trae los stocks de la API y los guarda en la DB (de a una página a la vez)
func StartFetchHandler(c *gin.Context) {
	// Genero ID único de tarea
	id := uuid.New().String()
	// Inicializo estructura de tarea en progreso
	info := &TaskInfo{Status: "in-progress", PagesFetched: 0}

	// Hago lock del mutex del mapa de tareas
	tasksMu.Lock()
	// Escribo la tarea en el mapa de tareas (se pasa por referencia)
	tasks[id] = info
	tasksMu.Unlock()

	// Inicio goroutine (se ejecuta en paralelo)
	go func() {
		// Marco, al último, el estado de la tarea
		defer func() {
			tasksMu.Lock()
			if info.Status != "error" {
				info.Status = "done"
			}
			tasksMu.Unlock()
		}()

		next := ""
		for batch := 0; ; batch++ {
			// Traigo la página actual (ya en formato de [] de Stock en model.go)
			stocks, nextPage, err := fetchPage(next)
			if err != nil {
				// Marco en caso de error
				tasksMu.Lock()
				info.Status = "error"
				info.Error = err.Error()
				tasksMu.Unlock()
				return
			}
			// Convierto los items traidos de la página a estructura Stock
			//stocks := convertItemsToStocks(r.Items)
			// Guardo los stocks en la BD
			if err := saveStocks(stocks); err != nil {
				// Marco en caso de error
				tasksMu.Lock()
				info.Status = "error"
				info.Error = err.Error()
				tasksMu.Unlock()
				return
			}

			// Actualizo progreso de la tarea
			tasksMu.Lock()
			info.PagesFetched = batch + 1
			tasksMu.Unlock()

			// Termino en caso de que no queden más páginas
			if nextPage == "" {
				break
			}
			next = nextPage
		}
	}()

	// Respondo al cliente con el id de la tarea lanzada
	c.JSON(http.StatusAccepted, gin.H{"task_id": id})
}

// Handler para consultar el estado de una tarea
func FetchStatusHandler(c *gin.Context) {
	// Obtengo id desde URL
	id := c.Param("id")
	// Mutex de lectura, leo la info (ok indica si encontró la tarea)
	tasksMu.RLock()
	info, ok := tasks[id]
	tasksMu.RUnlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}
	// Devuelvo info de la tarea encontrada
	c.JSON(http.StatusOK, info)
}

// Manejador para listar los stocks guardados en la BD
func listStocksHandler(c *gin.Context) {
	// Contador de la petición actual
	draw := c.DefaultQuery("draw", "1") //Obtengo parámetros y le asigno un default si no existen
	// Índice del primer registro a mostrar
	startStr := c.DefaultQuery("start", "0")
	// Cantidad de registros por pag
	lengthStr := c.DefaultQuery("length", "10")
	// Texto a buscar
	search := c.Query("search[value]") //Le asigno "" si no existe

	// Convierto a enteros los valores obtenidos
	start, _ := strconv.Atoi(startStr)
	length, _ := strconv.Atoi(lengthStr)

	// Obtengo datos para el ordenamiento de DataTables
	orderColumnIndex := c.Query("order[0][column]")                                // Índice de la col a ordenar
	orderDir := c.DefaultQuery("order[0][dir]", "asc")                             // "asc" o "desc"
	orderColumnName := c.Query(fmt.Sprintf("columns[%s][data]", orderColumnIndex)) // Nombre de la columna a ordenar
	if orderColumnName == "" {
		orderColumnName = "ticker" // Columna por defecto
	}

	// Cuento el total de datos en tabla de stocks
	var total int64
	db.Table("stocks").Count(&total)

	// Query base
	query := db.Table("stocks")

	// Si hay texto a buscar, filtro resultados (solo busco en "ticker" o en "company")
	if search != "" {
		query = query.Where("ticker ILIKE ? OR company ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Cuento la cantidad de resultados filtrados
	var filtered int64
	query.Count(&filtered)

	var stocks []Stock
	query.
		Order(fmt.Sprintf("%s %s", orderColumnName, orderDir)). // Aplico orden y dirección
		Offset(start).                                          // Indico de qué registro empezar
		Limit(length).                                          // Indico cuántos registros traer
		Find(&stocks)                                           // Ejecuto consulta y guardo en variable stocks

	// Devuelvo respuesta con el formato que requiere DataTables
	c.JSON(http.StatusOK, gin.H{
		"draw":            draw,
		"recordsTotal":    total,
		"recordsFiltered": filtered,
		"data":            stocks,
	})
}

func StockDetailHandler(c *gin.Context) {
    ticker := c.Param("ticker")

    // 2.1) parseo y validación de filtros
    filters, err := parseHistoryFilters(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 2.2) Traer el registro de Stock
    var stock Stock
    res := db.First(&stock, "ticker = ?", ticker)
    // Si no encontró ningún registro…
    if res.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "stock not found"})
        return
    }
    // Si hubo un error distinto…
    if res.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
        return
    }

    // 2.3) Traer histórico con filtros
    history, err := getHistory(ticker, filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 2.5) Calcular Risk/Reward
    rr := calcRiskReward(history)

    // 2.6) Distribución de ratings en todos los stocks
    ratingDist, err := getRatingDistribution()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 3) Armar y devolver la respuesta
    resp := StockDetailResponse{
        Stock:              stock,
        History:            history,
        RiskReward:         rr,
        RatingDistribution: ratingDist,
    }
    c.JSON(http.StatusOK, resp)
}


func parseHistoryFilters(c *gin.Context) (HistoryFilters, error) {
    days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

    var hf HistoryFilters
    hf.Days = days

    if sd := c.Query("start_date"); sd != "" {
        t, err := time.Parse("2006-01-02", sd)
        if err != nil {
            return hf, err
        }
        hf.StartDate = &t
    }
    if ed := c.Query("end_date"); ed != "" {
        t, err := time.Parse("2006-01-02", ed)
        if err != nil {
            return hf, err
        }
        hf.EndDate = &t
    }
    if mp := c.Query("min_price"); mp != "" {
        v, err := strconv.ParseFloat(mp, 64)
        if err != nil {
            return hf, err
        }
        hf.MinPrice = &v
    }
    if mp := c.Query("max_price"); mp != "" {
        v, err := strconv.ParseFloat(mp, 64)
        if err != nil {
            return hf, err
        }
        hf.MaxPrice = &v
    }
    if mv := c.Query("min_volume"); mv != "" {
        v, err := strconv.ParseInt(mv, 10, 64)
        if err != nil {
            return hf, err
        }
        hf.MinVolume = &v
    }
    hf.OrderDesc = c.DefaultQuery("order", "asc") == "desc"
    return hf, nil
}

func getHistory(ticker string, f HistoryFilters) ([]HistoricalPoint, error) {
    q := db.Model(&HistoricalPoint{}).
        Where("ticker = ?", ticker)

    if f.StartDate != nil {
        q = q.Where("date >= ?", f.StartDate)
    }
    if f.EndDate != nil {
        q = q.Where("date <= ?", f.EndDate)
    }
    if f.MinPrice != nil {
        q = q.Where("close >= ?", *f.MinPrice)
    }
    if f.MaxPrice != nil {
        q = q.Where("close <= ?", *f.MaxPrice)
    }
    if f.MinVolume != nil {
        q = q.Where("volume >= ?", *f.MinVolume)
    }

    orderDir := "asc"
    if f.OrderDesc {
        orderDir = "desc"
    }
    q = q.Order("date " + orderDir).
        Limit(f.Days)

    var pts []HistoricalPoint
		// Busco los historicos
    if err := q.Find(&pts).Error; err != nil {
        return nil, err
    }

    // Si pedimos desc, invertimos el slice para devolver asc por JSON
    if f.OrderDesc {
        for i, j := 0, len(pts)-1; i < j; i, j = i+1, j-1 {
            pts[i], pts[j] = pts[j], pts[i]
        }
    }
    return pts, nil
}

func calcRiskReward(history []HistoricalPoint) RiskRewardData {
    labels := make([]string, len(history))
    vols := make([]float64, len(history))
    pots := make([]float64, len(history))

    for i, pt := range history {
        labels[i] = pt.Date.Format("2006-01-02")
        vols[i] = (pt.High - pt.Low) / pt.Open * 100
        pots[i] = (pt.Close - pt.Open) / pt.Open * 100
    }
    return RiskRewardData{
        Labels:       labels,
        Volatilities: vols,
        Potentials:   pots,
    }
}

// Obtengo los ratings agrupados y contada su cantidad
func getRatingDistribution() (map[string]int, error) {
    rows, err := db.Model(&Stock{}).
        Select("rating_to, count(*) as cnt").
        Group("rating_to").
        Rows()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

		// Guardo en un map
    m := make(map[string]int)
    for rows.Next() {
        var rating string
        var cnt int
        rows.Scan(&rating, &cnt)
        m[rating] = cnt
    }
    return m, nil
}

// Inicio goroutine para traer los datos históricos de cada stock desde Yahoo Finance
func StartEnrichHandler(c *gin.Context) {
	// Genero un ID único para la tarea
	taskID := uuid.New().String()
	// Inicializo estructura de tarea en progreso
	info := &TaskInfo{Status: "in-progress", PagesFetched: 0}

	// Hago lock del mutex del mapa de tareas
	tasksMu.Lock()
	// Escribo la tarea en el mapa de tareas (se pasa por referencia)
	tasks[taskID] = info
	tasksMu.Unlock()

	// Inicio goroutine
	go func(taskID string) {
		// Recorro los stock, e intento traer los datos
		err := fetchAllHistories(taskID)
		tasksMu.Lock()
		defer tasksMu.Unlock()
		ti := tasks[taskID]
		if err != nil {
			ti.Status = "error"
			ti.Error = err.Error()
		} else {
			ti.Status = "done"
		}
	}(taskID)
	// Devuelvo de inmediato el ID
	c.JSON(202, gin.H{"task_id": taskID})
}

func RecalculateRecommendationsHandler(c *gin.Context) {
	// Disparo goroutine para recalcular las recomendaciones
    go func() {
        if err := RecalculateRecommendations(); err != nil {
            log.Printf("Error recalculando recomendaciones: %v", err)
        }
    }()
    c.JSON(http.StatusAccepted, gin.H{"status": "started"})
}

func TopRecommendationsHandler(c *gin.Context) {
    var recs []Recommendation
    if err := db.Order("score DESC").Limit(20).Find(&recs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, recs)
}
