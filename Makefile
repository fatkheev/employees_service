.PHONY: up down build run help

up: ## Запускает docker-compose
	@docker-compose up --build

down: ## Останавливает и удаляет контейнеры
	@docker-compose down

build: ## Собирает образ приложения
	@docker-compose build

run: ## Запускает приложение локально
	@go run main.go

help: ## Выводит помощь по использованию make
	@echo "Использование:"
	@echo "  make [command]"
	@echo ""
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help