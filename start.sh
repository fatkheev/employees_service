#!/bin/sh

# Применение миграций
echo "Applying migrations..."
/go/bin/migrate -path ./docs/migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" down
/go/bin/migrate -path ./docs/migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up

# # Запуск приложения
# echo "Starting application..."
# /employee-service
