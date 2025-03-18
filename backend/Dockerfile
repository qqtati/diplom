FROM golang:1.23-alpine

WORKDIR /app

# Копируем файлы с зависимостями
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Создаем директорию для загруженных файлов
RUN mkdir -p uploads/homework

# Собираем приложение
RUN go build -o main ./cmd/api

# Запускаем приложение
CMD ["./main"]
