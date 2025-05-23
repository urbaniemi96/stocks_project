package handlers

import (
	//"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/urbaniemi96/stocks_project/backend/db"
	"github.com/urbaniemi96/stocks_project/backend/detail"
	"github.com/urbaniemi96/stocks_project/backend/fetcher"
	"github.com/urbaniemi96/stocks_project/backend/model"
	"github.com/urbaniemi96/stocks_project/backend/recommender"
	"github.com/urbaniemi96/stocks_project/backend/tasks"
	"log"
	"net/http"
	"strconv"
)

// Inicia una goroutine que trae los stocks de la API y los guarda en la DB (de a una página a la vez)
func StartFetchHandler(c *gin.Context) {
	// Genero ID único de tarea
	id := uuid.New().String()
	// Inicializo estructura de tarea en progreso
	info := &tasks.TaskInfo{Status: "in-progress", PagesFetched: 0, Error: ""}

	// Hago lock del mutex del mapa de tareas
	tasks.TasksMu.Lock()
	// Escribo la tarea en el mapa de tareas (se pasa por referencia)
	tasks.Tasks[id] = info
	tasks.TasksMu.Unlock()

	// Inicio goroutine (se ejecuta en paralelo)
	go func() {
		// Marco, al último, el estado de la tarea
		defer func() {
			tasks.TasksMu.Lock()
			if info.Status != "error" {
				info.Status = "done"
			}
			tasks.TasksMu.Unlock()
		}()

		next := ""
		for batch := 0; ; batch++ {
			// Traigo la página actual (ya en formato de [] de Stock en model.go)
			stocks, nextPage, err := fetcher.FetchPage(next)
			if err != nil {
				// Marco en caso de error
				tasks.TasksMu.Lock()
				info.Status = "error"
				info.Error = err.Error()
				tasks.TasksMu.Unlock()
				return
			}
			// Guardo los stocks en la BD
			if err := db.SaveStocks(stocks); err != nil {
				// Marco en caso de error
				tasks.TasksMu.Lock()
				info.Status = "error"
				info.Error = err.Error()
				tasks.TasksMu.Unlock()
				return
			}

			// Actualizo progreso de la tarea
			tasks.TasksMu.Lock()
			info.PagesFetched = batch + 1
			tasks.TasksMu.Unlock()

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
	tasks.TasksMu.RLock()
	info, ok := tasks.Tasks[id]
	tasks.TasksMu.RUnlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}
	// Devuelvo info de la tarea encontrada
	c.JSON(http.StatusOK, info)
}

// Manejador para listar los stocks guardados en la BD
func ListStocksHandler(c *gin.Context) {
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
	db.DB.Table("stocks").Count(&total)

	// Query base
	query := db.DB.Table("stocks")

	// Si hay texto a buscar, filtro resultados (solo busco en "ticker" o en "company")
	if search != "" {
		query = query.Where("ticker ILIKE ? OR company ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Cuento la cantidad de resultados filtrados
	var filtered int64
	query.Count(&filtered)

	var stocks []model.Stock
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

// Traigo históricos de un ticker específico
func StockDetailHandler(c *gin.Context) {
	ticker := c.Param("ticker")

	// Validación de filtros
	filters, err := detail.ParseHistoryFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Traigo stock de ese ticker
	var stock model.Stock
	res := db.DB.First(&stock, "ticker = ?", ticker)
	// Si no encuentro
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "stock not found"})
		return
	}
	// Si hubo otro error
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	// Traigo el histórico con los filtros (los únicos funcionales los de fechas)
	history, err := detail.GetHistory(ticker, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculo riesgo / recompensa
	rr := detail.CalcRiskReward(history)

	// Distribución de ratings (en desuso)
	ratingDist, err := detail.GetRatingDistribution()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Estructura de la respuesta
	resp := detail.StockDetailResponse{
		Stock:              stock,
		History:            history,
		RiskReward:         rr,
		RatingDistribution: ratingDist,
	}
	c.JSON(http.StatusOK, resp)
}

// Inicio goroutine para traer los datos históricos de cada stock desde Yahoo Finance
func StartEnrichHandler(c *gin.Context) {
	// Genero un ID único para la tarea
	taskID := uuid.New().String()
	// Inicializo estructura de tarea en progreso
	info := &tasks.TaskInfo{Status: "in-progress", PagesFetched: 0, Error: ""}

	// Hago lock del mutex del mapa de tareas
	tasks.TasksMu.Lock()
	// Escribo la tarea en el mapa de tareas (se pasa por referencia)
	tasks.Tasks[taskID] = info
	tasks.TasksMu.Unlock()

	// Inicio goroutine
	go func(taskID string) {
		// Recorro los stock, e intento traer los datos
		err := fetcher.FetchAllHistories(taskID)
		tasks.TasksMu.Lock()
		defer tasks.TasksMu.Unlock()
		ti := tasks.Tasks[taskID]
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

// Handler para recalcular las recomendaciones de stocks
func RecalculateRecommendationsHandler(c *gin.Context) {
	// Disparo goroutine para recalcular las recomendaciones
	go func() {
		if err := recommender.RecalculateRecommendations(); err != nil {
			log.Printf("Error recalculando recomendaciones: %v", err)
		}
	}()
	c.JSON(http.StatusAccepted, gin.H{"status": "started"})
}

// Obtengo el top 20 recomendaciones
func TopRecommendationsHandler(c *gin.Context) {
	var recs []model.Recommendation
	if err := db.DB.Order("score DESC").Limit(20).Find(&recs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recs)
}
