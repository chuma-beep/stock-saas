# Stock SaaS - Backend API



![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-316192?style=flat&logo=postgresql)
![Status](https://img.shields.io/badge/Status-Live-success)

A Go-based REST API for comparing stock performance across different time periods with AI-powered analysis.

## Why this project exists
This project demonstrates how to build a production-ready Go backend
that combines external data sources with AI inference.

The focus is not model training, but:
- Reliable API design
- Data caching and cost control
- AI used as a system component
- Clean separation of concerns


## 🚀 Features

- Fetch historical stock data from Alpha Vantage
- Compare multiple stocks side-by-side
- Calculate percentage changes, volatility, and correlations
- AI-powered insights using Groq LLaMA 3.3
- PostgreSQL database for caching stock data
- RESTful API with JSON responses

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **Framework:** Gin
- **Database:** PostgreSQ
- **APIs:** Alpha Vantage (stock data), Groq (AI analysis)
- **Deployment:** Render.com


## Request Flow

```mermaid
flowchart LR
Client --> API
API --> PostgreSQL
API --> AlphaVantage
API --> Groq
PostgreSQL --> API
Groq --> API
API --> Client
```

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- Alpha Vantage API key (free tier: https://www.alphavantage.co/)
- Groq API key (free tier: https://console.groq.com/)

## 🔧 Installation

### 1. Clone the repository
```bash
git clone https://github.com/chuma-beep/stock-saas.git
cd stock-saas
```

### 2. Install dependencies
```bash
go mod download
```

### 3. Set up environment variables

Create a `.env` file in the root directory:
```env
DATABASE_URL=postgresql://username:password@localhost:5432/stocksaas?sslmode=disable
ALPHA_VANTAGE_API_KEY=your_alpha_vantage_key
GROQ_API_KEY=your_groq_key
PORT=8080
```

### 4. Set up the database
```bash
# Create database
createdb stocksaas

# Apply schema
psql stocksaas < scripts/schema.sql
```

### 5. Run the server
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

## 📡 API Endpoints

### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "message": "Stock SaaS API is running"
}
```

---

### Fetch Stock Data
```http
GET /fetch/:ticker
```

Fetches and stores historical data for a stock (last 100 trading days).

**Parameters:**
- `ticker` (path) - Stock symbol (e.g., AAPL, MSFT)

**Example:**
```bash
curl http://localhost:8080/fetch/AAPL
```

**Response:**
```json
{
  "message": "Stock data fetched and stored",
  "ticker": "AAPL",
  "records": 100
}
```

---

### Get Stock Data
```http
GET /stock?ticker={ticker}&start={start_date}&end={end_date}
```

Retrieves stored stock data for a specific period.

**Query Parameters:**
- `ticker` - Stock symbol
- `start` - Start date (YYYY-MM-DD)
- `end` - End date (YYYY-MM-DD)

**Example:**
```bash
curl "http://localhost:8080/stock?ticker=AAPL&start=2025-11-01&end=2025-12-11"
```

**Response:**
```json
{
  "ticker": "AAPL",
  "start_date": "2025-11-01",
  "end_date": "2025-12-11",
  "data": [...],
  "percent_change": 12.4
}
```


---

### Compare Stocks
```http
GET /compare?ticker1={ticker1}&ticker2={ticker2}&start={start_date}&end={end_date}
```

Compares two stocks over the same period.

**Query Parameters:**
- `ticker1` - First stock symbol
- `ticker2` - Second stock symbol
- `start` - Start date (YYYY-MM-DD)
- `end` - End date (YYYY-MM-DD)

**Example:**
```bash
curl "http://localhost:8080/compare?ticker1=AAPL&ticker2=MSFT&start=2025-11-01&end=2025-12-11"
```

**Response:**
```json
{
  "comparison": [
    {
      "ticker": "AAPL",
      "percent_change": 12.4,
      "data": [...]
    },
    {
      "ticker": "MSFT",
      "percent_change": 8.7,
      "data": [...]
    }
  ],
  "start_date": "2025-11-01",
  "end_date": "2025-12-11"
}
```

---

### Get Current Prices
```http
GET /current-prices
```

Returns the latest prices for popular stocks.

**Response:**
```json
{
  "stocks": [
    {
      "symbol": "AAPL",
      "price": 195.71,
      "change": 2.3
    },
    ...
  ]
}
```


### Why AI via API (Groq)
- Focus on AI integration, not model trai```

Generates AI-powered insights about a stock comparison.

**Request Body:**
```json
{
  "comparison": { ... },
  "preset": "Christmas Season"
}
```

**Response:**
```json
{
  "analysis": "• AAPL outperformed MSFT by 3.7% during this period...\n• Key drivers: Holiday sales momentum...\n• Correlation: 72% (highly correlated)...\n• Volatility: AAPL 15.3%, MSFT 12.1%..."
}
```


```
## Key Engineering Decisions

### Why Go
- Strong concurrency model
- Low-latency HTTP services
- Simple deployment (single binary)

### Why PostgreSQL instead of Redis
- Persistent caching
- Queryable historical ranges
- Simpler infra for early-stage systems

### Why AI via API (Groq)
- Focus on AI integration, not model training
- Lower operational complexity

---

## 📁 Project Structure
```
stock-saas/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── database/
│   │   └── db.go             # Database connection & queries
│   ├── handler/
│   │   └── analyze.go        # AI analysis handler
│   ├── handlers/
│   │   └── stock.go          # Stock data handlers
│   ├── models/
│   │   └── stock.go          # Data models
│   └── services/
│       └── alphavantage.go   # Alpha Vantage API client
├── scripts/
│   └── schema.sql            # Database schema
├── .env                      # Environment variables (not in git)
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Deployment

### Render.com (Recommended)

1. **Create PostgreSQL database** on Render
2. **Create Web Service** on Render
3. **Configure environment variables**:
   - `DATABASE_URL` (from Render PostgreSQL)
   - `ALPHA_VANTAGE_API_KEY`
   - `GROQ_API_KEY`
   - `PORT=8080`
4. **Build Command**: `go build -o bin/stock-saas cmd/api/main.go`
5. **Start Command**: `./bin/stock-saas`

Auto-deploys on git push!

## 💡 Usage Tips

### Caching Strategy
- Stock data is cached in PostgreSQL after first fetch
- Subsequent requests read from database (fast & free)
- Refetch only when you need updated data

### API Rate Limits
- **Alpha Vantage Free**: 25 requests/day, 5/minute
- **Groq Free**: 30 requests/minute
- Pre-fetch popular stocks to minimize API usage

### Performance
- First request after inactivity: ~30s (Render cold start)
- Subsequent requests: <500ms
- Use cron-job.org to keep service warm

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


## 👤 Author

**Chukwuma Wisdom Anwaegbu**
- GitHub: [@chuma-beep](https://github.com/chuma-beep)
- Email: wisboynelson123@gmail.com chukwumawisdomanwaegbu@gmail.com

## 🙏 Acknowledgments

- [Alpha Vantage](https://www.alphavantage.co/) for stock market data
- [Groq](https://groq.com/) for AI inference
- [Gin](https://gin-gonic.com/) for the Go web framework
