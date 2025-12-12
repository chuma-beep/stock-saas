# Stock SaaS API - AI Coding Guidelines

## Architecture Overview
This is a Go-based REST API for stock market data using Gin framework. The application fetches daily stock data from Alpha Vantage API and stores it in PostgreSQL for querying and comparison.

**Key Components:**
- `cmd/api/main.go`: Entry point with Gin router setup and CORS
- `internal/handlers/stock.go`: HTTP handlers for fetch, query, and compare endpoints
- `internal/services/alphavantage.go`: External API client for Alpha Vantage data fetching
- `internal/database/db.go`: PostgreSQL operations with upsert logic
- `internal/models/stock.go`: Data structures for stock records and API responses
- `scripts/schema.sql`: Database schema with unique constraint on (ticker, date)

**Data Flow:**
1. `/fetch/:ticker` → `services.FetchStockData()` → Alpha Vantage API → `database.SaveStock()` (with ON CONFLICT DO UPDATE)
2. `/stock` and `/compare` → `database.GetStockData()` → Calculate percent change in handlers

## Key Patterns
- **Environment Configuration**: Use `godotenv.Load()` for `.env` file loading. Required vars: `DATABASE_URL`, `ALPHA_VANTAGE_API_KEY`, `PORT` (default 8080)
- **Database Connection**: Singleton `database.DB` var, connect in `main()`, defer close
- **Error Handling**: Return JSON errors in handlers, log.Printf for debugging
- **Stock Data Parsing**: Alpha Vantage JSON structure parsed with `fmt.Sscanf` for OHLCV values
- **Percent Change Calculation**: `((lastClose - firstClose) / firstClose) * 100` in handlers
- **Unique Constraints**: DB upsert prevents duplicate (ticker, date) entries

## Developer Workflows
- **Run Locally**: `go run cmd/api/main.go` (requires PostgreSQL running)
- **Database Setup**: `psql -d stocksaas -f scripts/schema.sql` (create DB first)
- **Test Endpoints**: `curl http://localhost:8080/fetch/IBM` then `curl "http://localhost:8080/stock?ticker=IBM&start=2024-01-01&end=2024-12-01"`
- **Debug API Issues**: Check raw response logging in `services.FetchStockData()` for rate limits or invalid keys
- **Port Conflicts**: Use `lsof -i :8080` to find process, `kill -9 PID` to free port

## Integration Points
- **Alpha Vantage API**: Daily time series endpoint, requires API key, handles rate limiting
- **PostgreSQL**: Local DB connection, no migrations (schema.sql manual)
- **External Dependencies**: Gin for routing, lib/pq for DB driver, godotenv for config

## Conventions
- Internal packages only (no external imports from internal/)
- Handler functions take `*gin.Context`, return JSON responses
- Database functions use prepared statements with $ placeholders
- No tests currently - validate manually with curl
- Logs use `log.Printf` for server events, `fmt.Printf` for API debug