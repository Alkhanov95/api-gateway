# ===== build stage =====
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарь
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# ===== runtime stage =====
FROM alpine:3.20
WORKDIR /app

# Копируем бинарь и конфиг
COPY --from=builder /app/app /app/app
COPY --from=builder /app/config /app/config

# Открываем порт приложения
EXPOSE 3000

CMD ["/app/app"]