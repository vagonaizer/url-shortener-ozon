FROM golang:1.20-alpine AS builder

WORKDIR /app

# Копируем файлы модулей и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник из каталога с main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/url-shortener

FROM alpine:latest

WORKDIR /app

# Копируем бинарник из builder-образа
COPY --from=builder /app/url-shortener .

# Открываем порты REST и gRPC
EXPOSE 8080 50051

ENTRYPOINT ["./url-shortener"]