# Используем официальный образ Golang в качестве базового
FROM golang:1.24 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы приложения
COPY . .

# Копируем файл конфигурации в папку config
COPY ./config/config.yaml ./config/config.yaml
COPY ./migrations_postgres ./migrations_postgres

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

# Используем минимальный образ для запуска
FROM alpine:latest

# Устанавливаем необходимые библиотеки
RUN apk --no-cache add ca-certificates

# Копируем скомпилированное приложение из предыдущего этапа
COPY --from=builder /app/main .

# Копируем файл конфигурации в образ
COPY --from=builder /app/config/config.yaml ./config/config.yaml
COPY --from=builder /app/migrations_postgres ./migrations_postgres
# Указываем команду для запуска приложения
CMD ["./main"]

# Открываем порт, на котором будет работать приложение
EXPOSE 8080