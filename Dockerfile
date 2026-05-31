# ---- стадия сборки ----
FROM golang:1.25-alpine AS builder

# Устанавливаем необходимые системные пакеты
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Генерируем Templ-файлы (если не сгенерированы заранее)
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1001 && templ generate
# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/voidsounds ./cmd/server

# ---- финальная стадия ----
FROM alpine:latest

# Устанавливаем переменные времени и сертификаты
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Копируем бинарник из стадии сборки
COPY --from=builder /build/voidsounds /app/voidsounds

# Копируем миграции (если нужно выполнять их внутри контейнера)
COPY --from=builder /build/migrations /app/migrations

# Открываем порт
EXPOSE 8080

# Запуск
CMD ["/app/voidsounds"]