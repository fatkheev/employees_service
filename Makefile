.PHONY: up down build run help migrate-up migrate-down clear

up: ## Запускает docker-compose
	@docker-compose up --build

down: ## Останавливает и удаляет контейнеры
	@docker-compose down

build: ## Собирает образ приложения
	@docker-compose build

run: ## Запускает приложение локально
	@go run main.go

migrate-up:
	goose -dir ./db/migrations postgres "host=db user=user password=password dbname=postgres sslmode=disable" up

migrate-down:
	goose -dir ./db/migrations postgres "host=db user=user password=password dbname=postgres sslmode=disable" down

clear:
	goose -dir ./db/migrations postgres "host=db user=user password=password dbname=postgres sslmode=disable" up-to 2023053102

help: ## Выводит помощь по использованию make
	@echo "Использование:"
	@echo "  make [command]"
	@echo ""
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help