package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	var err error
	connStr := os.Getenv("DATABASE_URL")

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	for i := 0; i < 10; i++ {
		if err = DB.Ping(); err == nil {
			log.Println("✅ Database connected successfully")
			return nil
		}
		log.Printf("⏳ Database connection attempt %d failed, retrying...", i+1)
		time.Sleep(time.Second * 2)
	}

	return fmt.Errorf("error connecting to database after 10 attempts: %w", err)
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func SaveStock(ticker string, date time.Time, open, high, low, close float64, volume int64) error {
	query := `
        INSERT INTO stocks (ticker, date, open, high, low, close, volume)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (ticker, date) DO UPDATE
        SET open = $3, high = $4, low = $5, close = $6, volume = $7
    `

	_, err := DB.Exec(query, ticker, date, open, high, low, close, volume)
	return err
}

func GetStockData(ticker, startDate, endDate string) ([]map[string]interface{}, error) {
	query := `
        SELECT ticker, date, open, high, low, close, volume
        FROM stocks
        WHERE ticker = $1 AND date BETWEEN $2 AND $3
        ORDER BY date ASC
    `

	rows, err := DB.Query(query, ticker, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var ticker string
		var date time.Time
		var open, high, low, close float64
		var volume int64

		if err := rows.Scan(&ticker, &date, &open, &high, &low, &close, &volume); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"ticker": ticker,
			"date":   date.Format("2006-01-02"),
			"open":   open,
			"high":   high,
			"low":    low,
			"close":  close,
			"volume": volume,
		})
	}

	return results, nil
}
