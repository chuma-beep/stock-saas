package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type StockComparison struct {
	Ticker        string  `json:"ticker"`
	PercentChange float64 `json:"percent_change"`
	Data          []struct {
		Date   string  `json:"date"`
		Close  float64 `json:"close"`
		Volume int64   `json:"volume"`
	} `json:"data"`
}

type ComparisonResponse struct {
	StartDate  string            `json:"start_date"`
	EndDate    string            `json:"end_date"`
	Comparison []StockComparison `json:"comparison"`
}

type AnalyzeRequest struct {
	Comparison ComparisonResponse `json:"comparison"`
	Preset     string             `json:"preset"`
}

type AnalyzeResponse struct {
	Analysis string `json:"analysis"`
}

// Groq structs
type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Model       string        `json:"model"`
	Messages    []GroqMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type GroqChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type GroqResponse struct {
	Choices []GroqChoice `json:"choices"`
}

// Helper: Pearson correlation on returns
func calculateCorrelation(returns1, returns2 []float64) float64 {
	n := len(returns1)
	if n != len(returns2) || n < 2 {
		return 0
	}
	var sum1, sum2, sum1Sq, sum2Sq, pSum float64
	for i := 0; i < n; i++ {
		r1, r2 := returns1[i], returns2[i]
		sum1 += r1
		sum2 += r2
		sum1Sq += r1 * r1
		sum2Sq += r2 * r2
		pSum += r1 * r2
	}
	num := pSum - (sum1 * sum2 / float64(n))
	den := math.Sqrt((sum1Sq - math.Pow(sum1, 2)/float64(n)) * (sum2Sq - math.Pow(sum2, 2)/float64(n)))
	if den == 0 {
		return 0
	}
	return num / den
}

// Helper: ternary for string
func ternary(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

// Helper: Annualized volatility (%)
func calculateVolatility(prices []float64) float64 {
	n := len(prices)
	if n < 2 {
		return 0
	}
	returns := make([]float64, n-1)
	for i := 1; i < n; i++ {
		returns[i-1] = (prices[i] - prices[i-1]) / prices[i-1]
	}
	mean := 0.0
	for _, r := range returns {
		mean += r
	}
	mean /= float64(len(returns))
	variance := 0.0
	for _, r := range returns {
		variance += math.Pow(r-mean, 2)
	}
	variance /= float64(len(returns) - 1)
	return math.Sqrt(variance) * math.Sqrt(252) * 100
}

// AnalyzeComparison handler
func AnalyzeComparison(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	comp := req.Comparison
	if len(comp.Comparison) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expected two stocks for comparison"})
		return
	}

	stockA, stockB := comp.Comparison[0], comp.Comparison[1]

	// Extract closes for calcs
	closesA := make([]float64, len(stockA.Data))
	closesB := make([]float64, len(stockB.Data))
	for i, d := range stockA.Data {
		closesA[i] = d.Close
	}
	for i, d := range stockB.Data {
		closesB[i] = d.Close
	}

	// Returns for correlation
	returnsA := make([]float64, len(closesA)-1)
	returnsB := make([]float64, len(closesB)-1)
	for i := 1; i < len(closesA); i++ {
		returnsA[i-1] = (closesA[i] - closesA[i-1]) / closesA[i-1]
		returnsB[i-1] = (closesB[i] - closesB[i-1]) / closesB[i-1]
	}

	corr := calculateCorrelation(returnsA, returnsB)
	volA := calculateVolatility(closesA)
	volB := calculateVolatility(closesB)

	winner := stockA.Ticker
	margin := stockA.PercentChange - stockB.PercentChange
	if margin < 0 {
		winner = stockB.Ticker
		margin = -margin
	}

	// Build prompt— this is where the magic happens. Structured for consistent, engaging output.
	prompt := fmt.Sprintf(`You are an expert stock analyst explaining to savvy retail traders.

Compare %s vs %s over %s to %s (%s period).

Key stats:
- %s: %.2f%% (%s)
- %s: %.2f%% (%s)
- Winner: %s by %.2f%%
- Correlation: %.1f%%
- Volatility: %s %.1f%%, %s %.1f%%

Explain in bullet points:
• Who won and why (tie to stats)
• Key drivers (volume, highs/lows, volatility)
• Seasonal context (e.g., holiday buzz, Santa rally)
• Risks/opportunities (hedging if correlated, options if volatile)
• Fun takeaway

Under 300 words, emojis for punch.`,
		stockA.Ticker, stockB.Ticker, comp.StartDate, comp.EndDate, req.Preset,
		stockA.Ticker, stockA.PercentChange, ternary(stockA.PercentChange >= 0, "up", "down"),
		stockB.Ticker, stockB.PercentChange, ternary(stockB.PercentChange >= 0, "up", "down"),
		winner, margin,
		corr*100,
		stockA.Ticker, volA, stockB.Ticker, volB,
	)

	groqReq := GroqRequest{
		Model:       "llama-3.3-70b-versatile",
		Messages:    []GroqMessage{{Role: "user", Content: prompt}},
		Temperature: 0.7,
		MaxTokens:   500,
	}

	body, err := json.Marshal(groqReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}
	greq, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	greq.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))
	greq.Header.Set("Content-Type", "application/json")

	gresp, err := client.Do(greq)
	if err != nil || gresp.StatusCode != http.StatusOK {
		c.JSON(http.StatusOK, AnalyzeResponse{Analysis: "🧠 AI insights temporarily unavailable (rate limit or network issue)—try again soon!"})
		return
	}
	defer gresp.Body.Close()

	var groqResp GroqResponse
	if err := json.NewDecoder(gresp.Body).Decode(&groqResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI response error"})
		return
	}

	analysis := "No analysis generated"
	if len(groqResp.Choices) > 0 {
		analysis = groqResp.Choices[0].Message.Content
	}

	c.JSON(http.StatusOK, AnalyzeResponse{Analysis: analysis})
}
