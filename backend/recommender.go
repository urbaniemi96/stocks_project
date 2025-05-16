// Archivo con algoritmo para recomendar inversion
//
// REVISAR Y MEJORAR TODO ESTO (lo dejo para después)
//

package main

import (
"time"
"gorm.io/gorm/clause"
)

func RecalculateRecommendations() error {
    // 1) Traer todo el histórico almacenado
    var hist []HistoricalPoint
    if err := db.Find(&hist).Error; err != nil {
        return err
    }

    // 2) Agrupar por ticker
    byTicker := make(map[string][]HistoricalPoint)
    for _, pt := range hist {
        byTicker[pt.Ticker] = append(byTicker[pt.Ticker], pt)
    }

    // 3) Calcular score para cada ticker
    recs := make([]Recommendation, 0, len(byTicker))
    for ticker, series := range byTicker {
			if len(series) == 0 {
				// Si no hay datos históricos, saltamos este ticker
				continue
			}
			var sumR, sumV float64
			for _, pt := range series {
					ret := (pt.Close - pt.Open) / pt.Open
					vol := (pt.High - pt.Low) / pt.Open
					sumR += ret
					sumV += vol
			}
			N := float64(len(series))
			avgR := sumR / N
			avgV := sumV / N
			sharpe := avgR / (avgV + 1e-6)

			// peso por rating actual
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

    // 4) Upsert masivo
    for _, r := range recs {
        db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&r)
    }
    return nil
}
