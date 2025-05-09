// Archivo con algoritmo para recomendar inversion
//
// REVISAR Y MEJORAR TODO ESTO (lo dejo para después)
//

package main

import (
  //"context"
  //"log"
)

/*func recommendBest(ctx context.Context) (*Stock, error) {
  rows, err := db.Query(ctx, `SELECT ticker, company, target_from, target_to, action, brokerage, rating_from, rating_to, time FROM stocks`)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var best *Stock
  var maxDiff float64

  for rows.Next() {
    var s Stock
    if err := rows.Scan(&s.Ticker, &s.Company, &s.TargetFrom, &s.TargetTo,
      &s.Action, &s.Brokerage, &s.RatingFrom, &s.RatingTo, &s.Time); err != nil {
      log.Println("scan:", err)
      continue
    }
    diff := s.TargetTo - s.TargetFrom
    if best == nil || diff > maxDiff {
      best = &s
      maxDiff = diff
    }
  }
  return best, nil
}*/