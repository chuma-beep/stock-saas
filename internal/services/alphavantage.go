package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type AlphaVantageResponse struct {
	MetaData   map[string]string            `json:"Meta Data"`
	TimeSeries map[string]map[string]string `json:"Time Series (Daily)"`
}

type StockData struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

func FetchStockData(ticker string) ([]StockData, error) {
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ALPHA_VANTAGE_API_KEY not set")
	}

	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s",
		ticker,
		apiKey,
	)

	fmt.Printf("🔍 Fetching from URL: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// DEBUG: Print the raw response (first 500 chars)
	if len(body) > 500 {
		fmt.Printf("📥 Raw API Response: %s...\n", string(body)[:500])
	} else {
		fmt.Printf("📥 Raw API Response: %s\n", string(body))
	}

	var avResp AlphaVantageResponse
	if err := json.Unmarshal(body, &avResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(avResp.TimeSeries) == 0 {
		return nil, fmt.Errorf("no data returned for ticker %s - API might be rate limited or key invalid", ticker)
	}

	var stocks []StockData
	for dateStr, values := range avResp.TimeSeries {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		var open, high, low, close float64
		var volume int64

		fmt.Sscanf(values["1. open"], "%f", &open)
		fmt.Sscanf(values["2. high"], "%f", &high)
		fmt.Sscanf(values["3. low"], "%f", &low)
		fmt.Sscanf(values["4. close"], "%f", &close)
		fmt.Sscanf(values["5. volume"], "%d", &volume)

		stocks = append(stocks, StockData{
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}

	fmt.Printf("✅ Successfully parsed %d stock records\n", len(stocks))
	return stocks, nil
}
