package recommender

import (
	"github.com/urbaniemi96/stocks_project/backend/db"
	"github.com/urbaniemi96/stocks_project/backend/model"
	"github.com/urbaniemi96/stocks_project/backend/tasks"
	"gorm.io/gorm/clause"
	"time"
)

// Recalculo el score de los stocks
func RecalculateRecommendations(taskID string) error {
	// Traigo el histórico
	var hist []model.HistoricalPoint
	if err := db.DB.Find(&hist).Error; err != nil {
		return err
	}

	// Agrupo por ticker
	byTicker := make(map[string][]model.HistoricalPoint)
	for _, pt := range hist {
		byTicker[pt.Ticker] = append(byTicker[pt.Ticker], pt)
	}

	// Calculo score de cada ticker
	//recs := make([]model.Recommendation, 0, len(byTicker))
	for ticker, series := range byTicker {
		if len(series) == 0 {
			// Si no hay datos históricos, continúo
			continue
		}
		var sumR, sumV float64
		for _, pt := range series {
			ret := (pt.Close - pt.Open) / pt.Open // rendimiento diario
			vol := (pt.High - pt.Low) / pt.Open   // volatilidad diaria
			sumR += ret
			sumV += vol
		}
		N := float64(len(series))
		avgR := sumR / N               // rendimiento medio
		avgV := sumV / N               // volatilidad media
		sharpe := avgR / (avgV + 1e-6) // ratio simplificado de Sharpe

		// ponderación por rating actual
		var s model.Stock
		if err := db.DB.Select("rating_to").First(&s, "ticker = ?", ticker).Error; err != nil {
			continue
		}
		weightMap := map[string]float64{"sell": 0.8, "hold": 1.0, "buy": 1.2}
		weight, ok := weightMap[s.RatingTo]
		if !ok {
			weight = 1.0
		}

		score := sharpe * weight

		rec := model.Recommendation{
			Ticker:    ticker,
			Score:     score,
			UpdatedAt: time.Now(),
		}
		db.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&rec)
		// Actualizo el contador de tickers procesados
		tasks.TasksMu.Lock()
		tasks.Tasks[taskID].PagesFetched += 1
		tasks.TasksMu.Unlock()
	}

	// Upsert masivo
	/*for _, r := range recs {
		db.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&r)
	}*/
	return nil
}
