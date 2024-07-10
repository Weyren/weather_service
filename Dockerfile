# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем все зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Переключаемся в директорию, где лежит main.go
WORKDIR /app/cmd

# Собираем приложение
RUN go build -o /app/main .

# Финальная стадия
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный бинарник из предыдущей стадии
COPY --from=builder /app/main .

# Открываем порт 8080 для внешнего мира
EXPOSE 8080

# Команда для запуска исполняемого файла
CMD ["./main"]
