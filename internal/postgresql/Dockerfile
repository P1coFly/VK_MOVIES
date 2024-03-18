# /internal/postgresql

FROM postgres:latest

# Копируем файл инициализации в каталог docker-entrypoint-initdb.d
COPY init.sql /docker-entrypoint-initdb.d/
