package main

import (
	"sync"
)

// Estructura para guardar el progreso de la tarea
type TaskInfo struct {
	Status       string `json:"status"`        // "in-progress" / "done" / "error"
	PagesFetched int    `json:"pages_fetched"` // cantidad de pag procesadas
	Error        string `json:"error"`
}

var (
	// Mapa de tareas pasadas como referencia (para que al modificar una tarea, se refleje también en el mapa)
	tasks = make(map[string]*TaskInfo)
	// Semáforo Rmutex para permitir lecturas concurrentes en el mapa de tareas, pero solo una escritura a la vez
	tasksMu sync.RWMutex
)
