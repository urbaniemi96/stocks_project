package detail

import (
	"github.com/gin-gonic/gin"
	"github.com/urbaniemi96/stocks_project/backend/db"
	"github.com/urbaniemi96/stocks_project/backend/model"
	"strconv"
	"time"
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
	Stock              model.Stock             `json:"stock"`
	History            []model.HistoricalPoint `json:"history"`
	RiskReward         RiskRewardData          `json:"riskReward"`
	RatingDistribution map[string]int          `json:"ratingDistribution"`
}

type RiskRewardData struct {
	Labels       []string  `json:"labels"`
	Volatilities []float64 `json:"volatilities"`
	Potentials   []float64 `json:"potentials"`
}

func ParseHistoryFilters(c *gin.Context) (HistoryFilters, error) {
	// Cantidad de días a traer (90 por defecto)
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))

	var hf HistoryFilters
	hf.Days = days

	// Armo la estructura de filtros históricos
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
	// Si ambas fechas están presentes, y start > end, las intercambio
	if hf.StartDate != nil && hf.EndDate != nil && hf.StartDate.After(*hf.EndDate) {
		hf.StartDate, hf.EndDate = hf.EndDate, hf.StartDate
	}
	// En desuso por haber quitado los filtros
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

func GetHistory(ticker string, f HistoryFilters) ([]model.HistoricalPoint, error) {
	// Armo consulta
	q := db.DB.Model(&model.HistoricalPoint{}).
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
	q = q.Order("date " + orderDir) // .Limit(f.Days) Saco el límite de días (por haber quitado el filtro)

	var pts []model.HistoricalPoint
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

func CalcRiskReward(history []model.HistoricalPoint) RiskRewardData {
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

// Obtengo los ratings agrupados y contada su cantidad (en desuso)
func GetRatingDistribution() (map[string]int, error) {
	rows, err := db.DB.Model(&model.Stock{}).
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
