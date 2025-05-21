package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func parseHistoryFilters(c *gin.Context) (HistoryFilters, error) {
    days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

    var hf HistoryFilters
    hf.Days = days

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

func getHistory(ticker string, f HistoryFilters) ([]HistoricalPoint, error) {
    q := db.Model(&HistoricalPoint{}).
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
    q = q.Order("date " + orderDir).
        Limit(f.Days)

    var pts []HistoricalPoint
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

func calcRiskReward(history []HistoricalPoint) RiskRewardData {
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

// Obtengo los ratings agrupados y contada su cantidad
func getRatingDistribution() (map[string]int, error) {
    rows, err := db.Model(&Stock{}).
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