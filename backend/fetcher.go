package main

import (
  //"context"
  "encoding/json"
  "io"
  //"log"
  "net/http"
  "time"
  "strconv"
  "strings"
) 

// Estructura de la respuesta de la API
type apiResponse struct {
  Items    []struct {
    Ticker     string `json:"ticker"`
    Company    string `json:"company"`
    TargetFrom string `json:"target_from"`
    TargetTo   string `json:"target_to"`
    Action     string `json:"action"`
    Brokerage  string `json:"brokerage"`
    RatingFrom string `json:"rating_from"`
    RatingTo   string `json:"rating_to"`
    Time       string `json:"time"`
  } `json:"items"`
  NextPage string `json:"next_page"`
}

// Traigo todos los stocks de la API
func fetchAllStocks() ([]Stock, error) {
  // Inicio array de estructuras tipo Stock (definidas en model.go)
  var all []Stock
  next := ""

  // Creo cliente http para mandar peticiones (podría usar http.Get() también)
  //
  // REVISAR: poner timeouts para la petición?
  //
  client := &http.Client{}

  // Cargo URL y KEY (las obtengo en el config)
  var apiURL = getAPIURL()
  var apiKey = getAPIKEY()

  for {
    url := apiURL
    // Agrego parámetro next_page para obtener los datos siguientes
    if next != "" {
      url += "?next_page=" + next
    }
    // Armo la solicitud (cuerpo null porque no envío nada)
    req, _ := http.NewRequest("GET", url, nil)
    // Agrego header de autenticación (Bearer es un tipo de autenticación HTTP) 
    req.Header.Add("Authorization", "Bearer "+apiKey)
    // Envío petición
    resp, err := client.Do(req)
    if err != nil {
      return nil, err
    }
	  // Leo el cuerpo de la respuesta
	  body, _ := io.ReadAll(resp.Body)
	  // Cierro el stream de datos del cuerpo (lo hago con defer para asegurarme que se ejecute al final)
    defer resp.Body.Close()

    var r apiResponse
	  // Convierto el JSON del body en estructura de Go y lo guardo en r
    if err := json.Unmarshal(body, &r); err != nil {
      return nil, err
    }

	  // Recorro los items de la respuesta (ignoro el índice)
    for _, it := range r.Items {
      fFrom, _ := parseDollar(it.TargetFrom)
      fTo, _   := parseDollar(it.TargetTo)
	    // Convierto el string time devuelto por la api (RFC3339 con nanosegundos) en un objeto time de go
      t, _     := time.Parse(time.RFC3339Nano, it.Time)
	    // Agrego al array de Stock los datos traídos
      all = append(all, Stock{
        Ticker:     it.Ticker,
        Company:    it.Company,
        TargetFrom: fFrom,
        TargetTo:   fTo,
        Action:     it.Action,
        Brokerage:  it.Brokerage,
        RatingFrom: it.RatingFrom,
        RatingTo:   it.RatingTo,
        Time:       t,
      })
    }
	  // Controlo si no hay más datos para traer desde la API y paro el ciclo for
    //PARCHE PARA QUE NO TRAIGA INFINITOS (después lo arreglo, ahora solo quiero que funciones Dios)
    if r.NextPage == "LAMR" {
      break
    }
    next = r.NextPage
  }
  return all, nil
}

// Convierto valores del tipo "$4.60" en 4.60
func parseDollar(s string) (float64, error) {
  return strconv.ParseFloat(strings.Trim(s, "$"), 64)
}
