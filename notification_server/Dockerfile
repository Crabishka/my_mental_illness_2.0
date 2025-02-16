# Первый этап: сборка приложения
FROM golang:1.22-alpine AS builder

# Установка необходимых инструментов
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Установка swag и генерация документации
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN $GOPATH/bin/swag init -g ./cmd/server/main.go --parseDependency --parseInternal

# Сборка приложения
RUN go build -o main ./cmd/server/

# Второй этап: финальный образ
FROM alpine:latest

RUN apk --no-cache add tzdata

# Создание непривилегированного пользователя
RUN adduser -D -g '' appuser

WORKDIR /app

# Копирование бинарного файла и документации из builder
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./main"] 