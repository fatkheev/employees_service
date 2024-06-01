include .env
MIGRATE=migrate -path ./docs/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

.PHONY: up down build run help migrate-up migrate-down clear

up: up-db up-app ## Запускает docker-compose

up-db: ## Запускает docker-compose для базы данных
	@docker-compose up -d db

up-app: ## Ждет готовности базы данных и запускает остальные контейнеры
	@echo "Ожидание запуска базы данных..."
	@sleep 10 # Задержка в 10 секунд для ожидания готовности базы данных
	@docker-compose up --build swagger app

down: ## Останавливает и удаляет контейнеры
	@docker-compose down

build: ## Собирает образ приложения
	@docker-compose build

run: ## Запускает приложение локально
	@go run main.go

test: ## Запускает тесты
	@go test ./handlers -v

migrate-up:
	@echo "Применение миграций..."
	@$(MIGRATE) down
	@$(MIGRATE) up

migrate-down:
	@echo "Откат миграций..."
	@$(MIGRATE) down

migrate-force:
	@echo "Принудительный сброс миграций..."
	@$(MIGRATE) force $(version)

help: ## Выводит помощь по использованию make
	@echo "Использование:"
	@echo "  make [command]"
	@echo ""
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
