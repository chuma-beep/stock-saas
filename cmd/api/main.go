package main

import (
	"log"
	"os"

	"github.com/chuma-beep/stock-saas/internal/database"
	"github.com/chuma-beep/stock-saas/internal/handler"
	"github.com/chuma-beep/stock-saas/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Auto-create table with Unique Constraint

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS stocks (
        id SERIAL PRIMARY KEY,
        ticker VARCHAR(10) NOT NULL,
        date DATE NOT NULL,
        open NUMERIC NOT NULL,
        high NUMERIC NOT NULL,
        low NUMERIC NOT NULL,
        close NUMERIC NOT NULL,
        volume BIGINT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        UNIQUE (ticker, date)
    );`

	if _, err := database.DB.Exec(createTableQuery); err != nil {
		log.Fatal("Failed to create database table:", err)
	}
	log.Println("✅ Database table check passed")

	// Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.Default())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Stock SaaS API is running",
		})
	})

	r.POST("/api/analyze", handler.AnalyzeComparison)

	// Stock routes
	r.GET("/fetch/:ticker", handlers.FetchAndStoreStock)
	r.GET("/stock", handlers.GetStock)
	r.GET("/compare", handlers.CompareStocks)
	r.GET("/current-prices", handlers.GetCurrentPrices)

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Println("📊 Available endpoints:")
	log.Println("  GET /health")
	log.Println("  GET /fetch/:ticker")
	log.Println("  GET /stock?ticker=AAPL&start=2024-01-01&end=2024-12-01")
	log.Println("  GET /compare?ticker1=AAPL&ticker2=MSFT&start=2024-01-01&end=2024-12-01")

	r.Run(":" + port)
}
