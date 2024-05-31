DB_HOST := localhost
DB_PORT := 5432
DB_USER := user
DB_PASSWORD := password
DB_NAME := postgres

.PHONY: up down build run help migrate-up migrate-down clear

up: ## Запускает docker-compose
	@docker-compose up --build

down: ## Останавливает и удаляет контейнеры
	@docker-compose down

build: ## Собирает образ приложения
	@docker-compose build

run: ## Запускает приложение локально
	@go run main.go

DB_URL=postgres://user:password@localhost:5432/postgres?sslmode=disable

migrate-up:
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f create_tables.sql

migrate-down:
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f drop_tables.sql

help: ## Выводит помощь по использованию make
	@echo "Использование:"
	@echo "  make [command]"
	@echo ""
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help