# Stock SaaS API Documentation

Base URL: `https://your-backend.onrender.com`

## Authentication

Currently, the API is public and doesn't require authentication.

## Rate Limits

- Alpha Vantage: 25 requests/day (free tier)
- Groq AI: 30 requests/minute (free tier)

## Endpoints

[Include all the endpoints from the Backend README above]

## Error Handling

### Error Response Format
```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes

- `200` - Success
- `400` - Bad Request (missing parameters)
- `404` - Not Found (stock data not available)
- `500` - Internal Server Error

## Examples

### Fetch and Compare Workflow
```bash
# 1. Fetch AAPL data
curl -X GET "https://your-backend.onrender.com/fetch/AAPL"

# 2. Fetch MSFT data
curl -X GET "https://your-backend.onrender.com/fetch/MSFT"

# 3. Compare them
curl -X GET "https://your-backend.onrender.com/compare?ticker1=AAPL&ticker2=MSFT&start=2025-11-01&end=2025-12-11"
```

## Data Models

### Stock Data Point
```typescript
{
  ticker: string;
  date: string;      // ISO 8601 format
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
}
```

### Comparison Response
```typescript
{
  comparison: [
    {
      ticker: string;
      percent_change: number;
      data: StockDataPoint[];
    }
  ];
  start_date: string;
  end_date: string;
}
```
