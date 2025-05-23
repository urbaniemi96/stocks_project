package main

import (
"time"
"gorm.io/gorm/clause"
)

// Recalculo el score de los stocks
func RecalculateRecommendations() error {
    // Traigo el histórico
    var hist []HistoricalPoint
    if err := db.Find(&hist).Error; err != nil {
        return err
    }

    // Agrupo por ticker
    byTicker := make(map[string][]HistoricalPoint)
    for _, pt := range hist {
        byTicker[pt.Ticker] = append(byTicker[pt.Ticker], pt)
    }

    // Calculo score de cada ticker
    recs := make([]Recommendation, 0, len(byTicker))
    for ticker, series := range byTicker {
			if len(series) == 0 {
				// Si no hay datos históricos, continúo
				continue
			}
			var sumR, sumV float64
			for _, pt := range series {
					ret := (pt.Close - pt.Open) / pt.Open // rendimiento diario
					vol := (pt.High - pt.Low) / pt.Open // volatilidad diaria
					sumR += ret
					sumV += vol
			}
			N := float64(len(series))
			avgR := sumR / N // rendimiento medio
			avgV := sumV / N // volatilidad media
			sharpe := avgR / (avgV + 1e-6) // ratio simplificado de Sharpe

			// ponderación por rating actual
			var s Stock
			if err := db.Select("rating_to").First(&s, "ticker = ?", ticker).Error; err != nil {
					continue
			}
			weightMap := map[string]float64{"sell": 0.8, "hold": 1.0, "buy": 1.2}
			weight, ok := weightMap[s.RatingTo]
			if !ok {
					weight = 1.0
			}

			score := sharpe * weight

			recs = append(recs, Recommendation{
					Ticker:    ticker,
					Score:     score,
					UpdatedAt: time.Now(),
			})
    }

    // Upsert masivo
    for _, r := range recs {
        db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&r)
    }
    return nil
}
