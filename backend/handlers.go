package main

import (
  //"context"
  "net/http"
  "github.com/gin-gonic/gin"
  "strconv"
  "github.com/google/uuid"
  "fmt"
)
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

// Handler para consultar el estado de la tarea
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


// Manejador para obtener los stocks
/*func getStocksHandler(c *gin.Context) {
  // Traigo los stocks
  stocks, err := fetchAllStocks()

  if err != nil {
    c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
    return
  }

  // Guardo los stocks en la DB
  if err := saveStocks(context.Background(), stocks); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  // Devuelvo todo Ok
  c.JSON(http.StatusOK, stocks)
}*/

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
  orderColumnIndex := c.Query("order[0][column]") // Índice de la col a ordenar
	orderDir := c.DefaultQuery("order[0][dir]", "asc") // "asc" o "desc"
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
		Offset(start). // Indico de qué registro empezar
		Limit(length). // Indico cuántos registros traer
		Find(&stocks) // Ejecuto consulta y guardo en variable stocks

  // Devuelvo respuesta con el formato que requiere DataTables
  c.JSON(http.StatusOK, gin.H{
		"draw":            draw,
		"recordsTotal":    total,
		"recordsFiltered": filtered,
		"data":            stocks,
	})
}

// Manejador para listar los stocks guardados en la BD
/*func listStocksHandler(c *gin.Context) {
  rows, _ := db.Query(context.Background(), `SELECT ticker, company, target_from, target_to, action, brokerage, rating_from, rating_to, time FROM stocks`)
  defer rows.Close()

  var all []Stock
  for rows.Next() {
    var s Stock
    rows.Scan(&s.Ticker, &s.Company, &s.TargetFrom, &s.TargetTo,
      &s.Action, &s.Brokerage, &s.RatingFrom, &s.RatingTo, &s.Time)
    all = append(all, s)
  }
  c.JSON(http.StatusOK, all)
}*/

// Manejador para recomendar mejores stocks (REVISAR JUNTO CON EL ALGORITMO)
/*func recommendHandler(c *gin.Context) {
  best, err := recommendBest(context.Background())
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, best)
}*/