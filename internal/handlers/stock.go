package handlers

import (
	"log"
	"net/http"

	"github.com/chuma-beep/stock-saas/internal/database"
	"github.com/chuma-beep/stock-saas/internal/services"
	"github.com/gin-gonic/gin"
)

func FetchAndStoreStock(c *gin.Context) {
	ticker := c.Param("ticker")

	if ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticker is required"})
		return
	}

	log.Printf("Fetching data for %s...", ticker)

	stockData, err := services.FetchStockData(ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Storing %d records for %s...", len(stockData), ticker)

	saved := 0
	for _, data := range stockData {
		err := database.SaveStock(
			ticker,
			data.Date,
			data.Open,
			data.High,
			data.Low,
			data.Close,
			data.Volume,
		)
		if err != nil {
			log.Printf("Error saving data: %v", err)
			continue
		}
		saved++
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock data fetched and stored",
		"ticker":  ticker,
		"records": saved,
	})
}

func GetStock(c *gin.Context) {
	ticker := c.Query("ticker")
	startDate := c.Query("start")
	endDate := c.Query("end")

	if ticker == "" || startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ticker, start, and end query parameters are required",
		})
		return
	}

	data, err := database.GetStockData(ticker, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(data) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No data found. Try fetching it first using /fetch/:ticker",
		})
		return
	}

	// Calculate percent change
	firstClose := data[0]["close"].(float64)
	lastClose := data[len(data)-1]["close"].(float64)
	percentChange := ((lastClose - firstClose) / firstClose) * 100

	c.JSON(http.StatusOK, gin.H{
		"ticker":         ticker,
		"start_date":     startDate,
		"end_date":       endDate,
		"data":           data,
		"percent_change": percentChange,
	})
}

func CompareStocks(c *gin.Context) {
	ticker1 := c.Query("ticker1")
	ticker2 := c.Query("ticker2")
	startDate := c.Query("start")
	endDate := c.Query("end")

	if ticker1 == "" || ticker2 == "" || startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ticker1, ticker2, start, and end are required",
		})
		return
	}

	data1, err := database.GetStockData(ticker1, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data2, err := database.GetStockData(ticker2, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(data1) == 0 || len(data2) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Missing data. Fetch stocks first using /fetch/:ticker",
		})
		return
	}

	// Calculate percent changes
	calc := func(data []map[string]interface{}) float64 {
		first := data[0]["close"].(float64)
		last := data[len(data)-1]["close"].(float64)
		return ((last - first) / first) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"comparison": []gin.H{
			{
				"ticker":         ticker1,
				"percent_change": calc(data1),
				"data":           data1,
			},
			{
				"ticker":         ticker2,
				"percent_change": calc(data2),
				"data":           data2,
			},
		},
		"start_date": startDate,
		"end_date":   endDate,
	})
}
