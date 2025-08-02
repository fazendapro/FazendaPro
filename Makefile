.PHONY: help build run test clean docker-build docker-run docker-stop migrate migrate-docker

APP_NAME=fazendapro-api
DOCKER_IMAGE=fazendapro-api

help: ## Mostra esta ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Compila a aplicação
	go build -o bin/$(APP_NAME) main.go

run: ## Executa a aplicação localmente
	go run main.go

test: ## Executa os testes
	go test ./...

clean: ## Remove arquivos de build e limpa o ambiente
	rm -rf bin/
	go clean

docker-build: ## Constrói a imagem Docker
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Executa a aplicação com Docker Compose
	docker-compose up -d

docker-stop: ## Para a aplicação Docker Compose
	docker-compose down

docker-logs: ## Mostra logs dos containers
	docker-compose logs -f

migrate: ## Executa migrações localmente
	@chmod +x scripts/migrate.sh
	./scripts/migrate.sh up

migrate-docker: ## Executa migrações via Docker
	@chmod +x scripts/docker-migrate.sh
	./scripts/docker-migrate.sh up

dev-setup: ## Configura ambiente de desenvolvimento
	@echo "Configurando ambiente de desenvolvimento..."
	@if [ ! -f .env ]; then \
		cp env.development .env; \
		echo "Arquivo .env criado com configurações de desenvolvimento."; \
	fi
	@chmod +x scripts/*.sh

dev-start: dev-setup ## Inicia ambiente de desenvolvimento completo
	@echo "Iniciando ambiente de desenvolvimento..."
	./scripts/dev.sh start

dev-quick: ## Inicia apenas a aplicação (banco já deve estar rodando)
	@echo "Iniciando aplicação..."
	go run main.go

logs: ## Mostra logs da aplicação
	docker-compose logs -f app

db-connect: ## Conecta ao banco de dados
	docker-compose exec postgres psql -U fazendapro_user -d fazendapro

db-reset: ## Reseta o banco de dados
	docker-compose down -v
	docker-compose up -d postgres
	@sleep 10
	./scripts/docker-migrate.sh up

prod-build: ## Constrói para produção
	docker build -t $(DOCKER_IMAGE):latest .

prod-deploy: ## Deploy em produção (exemplo)
	@echo "Deploy em produção - implementar conforme necessidade" 