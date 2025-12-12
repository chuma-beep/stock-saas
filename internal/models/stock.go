package models

import "time"

type Stock struct {
	ID        int       `json:"id"`
	Ticker    string    `json:"ticker"`
	Date      time.Time `json:"date"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int64     `json:"volume"`
	CreatedAt time.Time `json:"created_at"`
}

type StockResponse struct {
	Ticker        string  `json:"ticker"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	Data          []Stock `json:"data"`
	PercentChange float64 `json:"percent_change"`
}
