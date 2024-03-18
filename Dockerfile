# Используем официальный образ Golang в качестве базового образа
FROM golang:latest

# Установка переменной окружения для указания рабочей директории внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта внутрь контейнера
COPY . .

# Устанавливаем переменную окружения для конфигурационного файла
ENV CONFIG_PATH=./config/config.yml

# Собираем приложение
RUN go build -o app ./cmd/apiserver/main.go

# Устанавливаем порт, который будет открыт в контейнере
EXPOSE 8080

# Команда для запуска приложения
CMD ["./app"]