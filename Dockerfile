# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o weather_service ./cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/weather_service /app/weather_service
COPY config.yml /app/config.yml
COPY migrate.sql /app/migrate.sql
COPY static /app/static
EXPOSE 8080
CMD ["./weather_service"]
