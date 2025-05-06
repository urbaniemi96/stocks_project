package main

import (
  //"context"
  "encoding/json"
  //"fmt"
  //"io/ioutil"
  //"log"
  "net/http"
)
// Cargo URL y KEY (las obtengo en el config)
var apiURL = getAPIURL()
var apiKey = getAPIKEY() 

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

// Traigo todas las acciones de la API
func fetchAllStocks() ([]Stock, error) {
  var all []Stock
  next := ""
  client := &http.Client{}

  for {
    url := apiURL
    if next != "" {
      url += "?next_page=" + next
    }
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Add("Authorization", "Bearer "+apiKey)
    resp, err := client.Do(req)
    if err != nil {
      return nil, err
    }
    body, _ := ioutil.ReadAll(resp.Body)
    resp.Body.Close()

    var r apiResponse
    if err := json.Unmarshal(body, &r); err != nil {
      return nil, err
    }

    for _, it := range r.Items {
      fFrom, _ := parseDollar(it.TargetFrom)
      fTo, _   := parseDollar(it.TargetTo)
      t, _     := time.Parse(time.RFC3339Nano, it.Time)
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

    if r.NextPage == "" {
      break
    }
    next = r.NextPage
  }
  return all, nil
}

// parseDollar convierte "$4.20" → 4.20
func parseDollar(s string) (float64, error) {
  return strconv.ParseFloat(strings.Trim(s, "$"), 64)
}
