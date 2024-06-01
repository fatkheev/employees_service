#!/bin/sh
set -e

echo "Waiting for database to be ready..."
sleep 10  # Задержка в 10 секунд для ожидания готовности базы данных

echo "Starting application..."
exec /employee-service
