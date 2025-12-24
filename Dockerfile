# Stage 1: Builder – Use Go 1.25.5
FROM golang:1.25.5 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o stock-saas cmd/api/main.go

# Stage 2: Runtime – Tiny alpine
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/stock-saas .
EXPOSE 8080
CMD ["./stock-saas"]
