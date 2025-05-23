package fetcher

import (
	//"context"
	"encoding/json"
	"fmt"
	"github.com/urbaniemi96/stocks_project/backend/config"
	"github.com/urbaniemi96/stocks_project/backend/db"
	"github.com/urbaniemi96/stocks_project/backend/model"
	"github.com/urbaniemi96/stocks_project/backend/tasks"
	"gorm.io/gorm/clause"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	maxRetries     = 10              // cuántas veces reintentar como máximo
	initialBackoff = 2 * time.Second // backoff inicial
)

// Estructura de la respuesta de la API
type apiResponse struct {
	Items []struct {
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

// Traigo los datos de una página (nextPage) y los transformo en []Stock
func FetchPage(nextPage string) ([]model.Stock, string, error) {
	var r apiResponse
	// Inicio array de estructuras tipo Stock (definidas en model.go)
	var all []model.Stock
	url := config.GetAPIURL()
	key := config.GetAPIKEY()
	// Creo cliente http para mandar peticiones (podría usar http.Get() también)
	client := &http.Client{}
	if nextPage != "" {
		url += "?next_page=" + nextPage
	}
	// Armo la solicitud (cuerpo null porque no envío nada)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("new request: %w", err)
	}
	// Agrego header de autenticación (Bearer es un tipo de autenticación HTTP)
	req.Header.Add("Authorization", "Bearer "+key)
	// Envío petición
	resp, err := client.Do(req)
	if err != nil {
		return all, "", fmt.Errorf("http request: %w", err)
	}
	// Cierro el stream de datos del cuerpo (lo hago con defer para asegurarme que se ejecute al final)
	defer resp.Body.Close()
	// Leo el body de la respuesta
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return all, "", fmt.Errorf("read body: %w", err)
	}
	// Convierto el JSON del body en estructura de Go y lo guardo en r
	if err := json.Unmarshal(body, &r); err != nil {
		return all, "", fmt.Errorf("unmarshal JSON: %w", err)
	}
	// Recorro los r.items de la respuesta (ignoro el índice)
	for _, it := range r.Items {
		fFrom, err := parseDollar(it.TargetFrom)
		if err != nil {
			return all, "", fmt.Errorf("parse TargetFrom %q: %w", it.TargetFrom, err)
		}
		fTo, err := parseDollar(it.TargetTo)
		if err != nil {
			return all, "", fmt.Errorf("parse TargetTo %q: %w", it.TargetTo, err)
		}
		// Convierto el string time devuelto por la api (RFC3339 con nanosegundos) en un objeto time de go
		t, err := time.Parse(time.RFC3339Nano, it.Time)
		if err != nil {
			return all, "", fmt.Errorf("parse time %q: %w", it.Time, err)
		}
		// Agrego al array de Stock los datos traídos
		all = append(all, model.Stock{
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
	return all, r.NextPage, nil
}

// fetchHistory llama a Yahoo Finance para un ticker,
// trayendo datos de los últimos 90 días con intervalo diario.
func fetchHistory(ticker string) ([]model.HistoricalPoint, error) {
	//API de datos
	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v8/finance/chart/%s?range=3mo&interval=1d",
		ticker,
	)

	var resp *http.Response
	backoff := initialBackoff

	// intento con reintentos
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Armo request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error construyendo request para %s: %w", ticker, err)
		}

		// Seteo las cookies para yahoo finance
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Referer", "https://finance.yahoo.com/quote/"+ticker)

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error HTTP para %s: %w", ticker, err)
		}

		// si nos dan 200 OK, seguimos al parsing
		if resp.StatusCode == http.StatusOK {
			break
		}

		// Si es 429, nos toca backoff manual
		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close() // siempre cerrar el body antes de dormir
			log.Printf("HTTP 429 para %s, intentando de nuevo en %v…", ticker, backoff)
			time.Sleep(backoff)
			backoff *= 2 // backoff exponencial
			continue
		}

		// otros errores HTTP los devolvemos directo con cuerpo
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("error HTTP %d para %s: %s",
			resp.StatusCode, ticker, string(body))
	}

	// Si no devuelve nada, se avisa
	if resp == nil {
		return nil, fmt.Errorf("no se obtuvo JSON válido para %s tras %d intentos", ticker, maxRetries)
	}
	defer resp.Body.Close()

	// Si el status code es distinto de Ok
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		//fmt.Println(resp)
		return nil, fmt.Errorf("error HTTP %d para %s: %s", resp.StatusCode, ticker, string(bodyBytes))
	}

	// Estructura parcial del JSON que nos interesa
	var result struct {
		Chart struct {
			Result []struct {
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Close  []float64 `json:"close"`
						Open   []float64 `json:"open"`
						High   []float64 `json:"high"`
						Low    []float64 `json:"low"`
						Volume []int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			} `json:"result"`
			Error interface{} `json:"error"`
		} `json:"chart"`
	}

	// Decodifico la respuesta de la API y la guardo en la estructura
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parseando JSON para %s: %w", ticker, err)
	}

	// Si no hay datos para el ticker
	if len(result.Chart.Result) == 0 {
		return nil, fmt.Errorf("sin datos para %s", ticker)
	}

	// Primera serie de datos (porque busco de a 1 ticker)
	serie := result.Chart.Result[0]
	// Slice de modelo HistoricalPoint, del tamaño de la serie
	points := make([]model.HistoricalPoint, len(serie.Timestamp))
	// Yahoo devuelve los datos en series paralelas "timestamp" e "indicators", por eso recorro los timestamps
	for i, ts := range serie.Timestamp {
		points[i] = model.HistoricalPoint{
			Ticker: ticker,
			// Transformo el ts en objeto time
			Date: time.Unix(ts, 0),
			// Obtengo los indicadores por el índice
			Close:  serie.Indicators.Quote[0].Close[i],
			Open:   serie.Indicators.Quote[0].Open[i],
			High:   serie.Indicators.Quote[0].High[i],
			Low:    serie.Indicators.Quote[0].Low[i],
			Volume: serie.Indicators.Quote[0].Volume[i],
		}
	}

	return points, nil
}

// fetchAllHistories recorre todos los tickers en DB,
// llama a fetchHistory con throttle de 1 cada 2s,
// guarda en price_histories, y continúa ante errores.
func FetchAllHistories(taskID string) error {
	//Traigo los stocks
	var stocks []model.Stock
	if err := db.DB.Find(&stocks).Error; err != nil {
		log.Fatalf("no pude leer tickers: %v", err)
		return err
	}

	// Recorro tickers
	for i, s := range stocks {
		// Traigo históricos
		points, err := fetchHistory(s.Ticker)
		if err != nil {
			log.Printf("ERROR al obtener datos para %s: %v", s.Ticker, err)
			time.Sleep(2 * time.Second)
			continue
		} else {
			// Guardo en db
			for _, pt := range points {
				if err := db.DB.Clauses(
					// Upsert
					clause.OnConflict{UpdateAll: true},
				).Create(&pt).Error; err != nil {
					log.Printf("ERROR saving %s [%s] Close=%.2f Open=%.2f High=%.2f Low=%.2f Vol=%d: %v\n",
						s.Ticker,
						pt.Date.Format("2006-01-02"),
						pt.Close, pt.Open, pt.High, pt.Low, pt.Volume,
						err,
					)
				}
			}
			log.Printf("Guardados %d puntos para %s\n", len(points), s.Ticker)
		}
		// Actualizo el contador de páginas (o tickers) procesados
		tasks.TasksMu.Lock()
		tasks.Tasks[taskID].PagesFetched = i + 1
		tasks.TasksMu.Unlock()
		// Espero 1 segundos antes de la siguiente petición
		time.Sleep(1 * time.Second)
	}

	log.Println("Proceso completo de fetchAllHistories")
	return nil
}

// Convierto valores del tipo "$4.60" en 4.60
func parseDollar(s string) (float64, error) {
	// Saco espacios
	clean := strings.TrimSpace(s)
	// Quito $
	clean = strings.TrimPrefix(clean, "$")
	// Quito separador de miles
	clean = strings.ReplaceAll(clean, ",", "")
	// Convierto a float64
	return strconv.ParseFloat(clean, 64)
}
