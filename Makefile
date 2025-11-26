# Makefile para FazendaPro API

.PHONY: help test test-coverage test-unit test-handlers clean install-deps run build migrate-docker db-reset swagger

# Variáveis
GO_VERSION := 1.24.2
COVERAGE_DIR := coverage
TEST_DIR := tests

# Ajuda
help:
	@echo "FazendaPro API - Comandos disponíveis:"
	@echo ""
	@echo "📦 Dependências:"
	@echo "  install-deps    - Instalar dependências do Go"
	@echo ""
	@echo "📚 Documentação:"
	@echo "  swagger         - Gerar documentação Swagger"
	@echo ""
	@echo "🧪 Testes:"
	@echo "  test            - Executar todos os testes"
	@echo "  test-coverage   - Executar testes com coverage"
	@echo "  test-unit       - Executar testes unitários"
	@echo "  test-handlers    - Executar testes dos handlers"
	@echo "  coverage-100    - Análise completa para 100% coverage"
	@echo "  generate-tests  - Gerar testes automaticamente"
	@echo ""
	@echo "🔧 Desenvolvimento:"
	@echo "  dev             - Inicia ambiente de desenvolvimento completo"
	@echo "  prod            - Inicia aplicação em produção"
	@echo "  quick           - Inicia apenas a aplicação (banco já deve estar rodando)"
	@echo "  run             - Executar a aplicação"
	@echo "  build           - Compilar a aplicação"
	@echo "  clean           - Limpar arquivos temporários"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  logs            - Mostra logs da aplicação"
	@echo "  db-connect      - Conecta ao banco de dados"
	@echo "  db-reset        - Recria o banco de dados do zero"
	@echo "  migrate-docker  - Executa migrações no ambiente Docker"
	@echo ""
	@echo "🚀 Produção:"
	@echo "  prod-build      - Constrói para produção"
	@echo "  prod-deploy     - Deploy em produção"
	@echo ""

# Instalar dependências
install-deps:
	@echo "📦 Instalando dependências..."
	go mod download
	go mod tidy
	@echo "✅ Dependências instaladas!"

# Executar todos os testes
test:
	@echo "🧪 Executando todos os testes..."
	go test -v ./...

# Executar testes com coverage
test-coverage:
	@echo "📊 Executando testes com coverage..."
	@mkdir -p $(COVERAGE_DIR)
	go test -v -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic -coverpkg=./internal/...,./tests/... ./$(TEST_DIR)/...
	go tool cover -func=$(COVERAGE_DIR)/coverage.out

# Executar testes unitários
test-unit:
	@echo "🔬 Executando testes unitários..."
	go test -v -short ./$(TEST_DIR)/...

# Executar testes dos handlers
test-handlers:
	@echo "🔐 Testando AuthHandler..."
	go test -v ./$(TEST_DIR)/handlers -run TestAuthHandler
	@echo ""
	@echo "👤 Testando UserHandler..."
	go test -v ./$(TEST_DIR)/handlers -run TestUserHandler
	@echo ""
	@echo "🐄 Testando AnimalHandler..."
	go test -v ./$(TEST_DIR)/handlers -run TestAnimalHandler
	@echo ""
	@echo "🥛 Testando MilkCollectionHandler..."
	go test -v ./$(TEST_DIR)/handlers -run TestMilkCollectionHandler
	@echo ""
	@echo "🔄 Testando ReproductionHandler..."
	go test -v ./$(TEST_DIR)/handlers -run TestReproductionHandler

# Executar a aplicação
run:
	@echo "🚀 Executando a aplicação..."
	go run main.go

# Compilar a aplicação
build:
	@echo "🔨 Compilando a aplicação..."
	go build -o fazendapro-api main.go
	@echo "✅ Aplicação compilada: fazendapro-api"

# Limpar arquivos temporários
clean:
	@echo "🧹 Limpando arquivos temporários..."
	rm -rf $(COVERAGE_DIR)
	rm -f fazendapro-api
	go clean
	@echo "✅ Limpeza concluída!"

# Executar script de coverage
coverage-script:
	@echo "📊 Executando script de coverage..."
	./scripts/test-coverage.sh

# Verificar qualidade do código
lint:
	@echo "🔍 Verificando qualidade do código..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint não instalado. Instale com: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Formatar código
fmt:
	@echo "🎨 Formatando código..."
	go fmt ./...
	@echo "✅ Código formatado!"

# Verificar dependências
deps-check:
	@echo "🔍 Verificando dependências..."
	go mod verify
	go list -u -m all
	@echo "✅ Dependências verificadas!"

# Instalar ferramentas de desenvolvimento
install-tools:
	@echo "🛠️  Instalando ferramentas de desenvolvimento..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/stretchr/testify@latest
	@echo "✅ Ferramentas instaladas!"

# Executar todos os checks
check: fmt lint test-coverage
	@echo "✅ Todos os checks executados com sucesso!"

# CI/CD pipeline
ci: install-deps fmt lint test-coverage
	@echo "✅ Pipeline CI/CD executado com sucesso!"

# =============================================================================
# COMANDOS DE DESENVOLVIMENTO
# =============================================================================

# Dar permissão de execução aos scripts
dev-setup:
	@echo "🔧 Configurando ambiente de desenvolvimento..."
	@chmod +x scripts/*.sh
	@echo "✅ Scripts configurados!"

# Inicia ambiente de desenvolvimento completo
dev: dev-setup ## Inicia ambiente de desenvolvimento completo
	@echo "Iniciando ambiente de desenvolvimento..."
	./scripts/dev.sh start
	ENV=development go run main.go

# Inicia aplicação em produção
prod: ## Inicia aplicação em produção
	@echo "Iniciando aplicação em produção..."
	cp env.production .env
	ENV=production go run main.go

# Inicia apenas a aplicação (banco já deve estar rodando)
quick: ## Inicia apenas a aplicação
	@echo "Iniciando aplicação..."
	go run main.go

# Mostra logs da aplicação
logs: ## Mostra logs da aplicação
	docker-compose logs -f app

# Conecta ao banco de dados
db-connect: ## Conecta ao banco de dados
	docker-compose exec postgres psql -U fazendapro_user -d fazendapro

# Recriar banco de dados (remove volumes e recria)
db-reset: ## Recria o banco de dados do zero
	@echo "🔄 Recriando banco de dados..."
	docker-compose down -v
	docker-compose up -d postgres
	@echo "⏳ Aguardando banco de dados inicializar..."
	@echo "   (aguardando PostgreSQL criar usuário e banco...)"
	@echo "   Isso pode levar até 60 segundos na primeira inicialização..."
	@timeout=120; \
	attempt=0; \
	internal_ready=0; \
	while [ $$timeout -gt 0 ]; do \
		attempt=$$((attempt + 1)); \
		if docker-compose exec -T postgres psql -U fazendapro_user -d fazendapro -c "SELECT 1;" > /dev/null 2>&1; then \
			internal_ready=1; \
			echo "✅ Banco de dados pronto e usuário criado! (tentativa $$attempt)"; \
			break; \
		fi; \
		if [ $$((attempt % 5)) -eq 0 ]; then \
			echo "   Aguardando... ($$timeout segundos restantes)"; \
		fi; \
		sleep 2; \
		timeout=$$((timeout - 2)); \
	done; \
	if [ $$internal_ready -eq 0 ]; then \
		echo "❌ Timeout aguardando banco de dados"; \
		echo "   Verifique os logs: docker-compose logs postgres"; \
		exit 1; \
	fi; \
	echo "   Aguardando healthcheck estar disponível..."; \
	timeout2=60; \
	healthcheck_ok=0; \
	while [ $$timeout2 -gt 0 ]; do \
		if docker-compose ps postgres 2>/dev/null | grep -q "healthy"; then \
			healthcheck_ok=1; \
			echo "✅ Healthcheck OK!"; \
			break; \
		fi; \
		sleep 2; \
		timeout2=$$((timeout2 - 2)); \
	done; \
	if [ $$healthcheck_ok -eq 0 ]; then \
		echo "⚠️  Healthcheck não ficou healthy, mas continuando..."; \
	fi; \
	echo "   Aguardando conexão externa estar disponível (pode levar alguns segundos)..."; \
	sleep 10
	@echo "📦 Executando migrações..."
	$(MAKE) migrate-docker
	@echo "✅ Banco de dados recriado com sucesso!"

# Executar migrações no Docker
migrate-docker: ## Executa migrações no ambiente Docker (usando rede Docker)
	@echo "📦 Executando migrações via Docker..."
	@if docker-compose ps postgres | grep -q "Up"; then \
		echo "   Usando serviço migration do docker-compose (rede Docker)..."; \
		docker-compose --profile migrate run --rm migration; \
	else \
		echo "⚠️  Container PostgreSQL não está rodando"; \
		exit 1; \
	fi

# Constrói para produção
prod-build: ## Constrói para produção
	docker build -t $(DOCKER_IMAGE):latest .

# Deploy em produção (exemplo)
prod-deploy: ## Deploy em produção
	@echo "Deploy em produção - implementar conforme necessidade"

# =============================================================================
# COMANDOS PARA 100% COVERAGE
# =============================================================================

# Análise completa para 100% coverage
coverage-100: ## Análise completa para alcançar 100% coverage
	@echo "🚀 Executando análise completa de coverage..."
	@chmod +x scripts/coverage-100.sh
	./scripts/coverage-100.sh

# Gerar testes automaticamente
generate-tests: ## Gerar testes automaticamente
	@echo "🤖 Gerando testes automaticamente..."
	@chmod +x scripts/generate-tests.sh
	./scripts/generate-tests.sh

# Executar testes gerados
test-generated: ## Executar testes gerados automaticamente
	@echo "🧪 Executando testes gerados..."
	@chmod +x scripts/run-generated-tests.sh
	./scripts/run-generated-tests.sh

# Análise de coverage avançada
coverage-analysis: ## Análise avançada de coverage
	@echo "📊 Executando análise avançada de coverage..."
	@mkdir -p coverage/analysis
	go test -coverprofile=coverage/analysis/coverage.out -covermode=atomic ./...
	go tool cover -html=coverage/analysis/coverage.out -o coverage/analysis/coverage.html
	go tool cover -func=coverage/analysis/coverage.out > coverage/analysis/coverage-func.txt
	@echo "📁 Relatórios gerados em: coverage/analysis/"

# Coverage com richgo (se disponível)
test-rich: ## Executar testes com richgo
	@echo "🎨 Executando testes com richgo..."
	@if command -v richgo >/dev/null 2>&1; then \
		richgo test -v -coverprofile=coverage/rich-coverage.out -covermode=atomic ./...; \
		go tool cover -html=coverage/rich-coverage.out -o coverage/rich-coverage.html; \
		echo "📁 Relatório richgo: coverage/rich-coverage.html"; \
	else \
		echo "⚠️  richgo não instalado. Instalando..."; \
		go install github.com/kyoh86/richgo@latest; \
		richgo test -v -coverprofile=coverage/rich-coverage.out -covermode=atomic ./...; \
	fi

# Pipeline completo para 100% coverage
coverage-pipeline: generate-tests coverage-100 test-generated ## Pipeline completo para 100% coverage
	@echo "✅ Pipeline de coverage executado com sucesso!"
	@echo "📊 Verifique os relatórios em: coverage/"

# =============================================================================
# COMANDOS DE DOCUMENTAÇÃO SWAGGER
# =============================================================================

# Gerar documentação Swagger
swagger: ## Gerar documentação Swagger
	@echo "📚 Gerando documentação Swagger..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g main.go -o docs --parseDependency --parseInternal; \
		echo "✅ Documentação Swagger gerada com sucesso!"; \
		echo "📖 Acesse: http://localhost:8080/swagger/index.html"; \
	elif [ -f "$(HOME)/go/bin/swag" ]; then \
		$(HOME)/go/bin/swag init -g main.go -o docs --parseDependency --parseInternal; \
		echo "✅ Documentação Swagger gerada com sucesso!"; \
		echo "📖 Acesse: http://localhost:8080/swagger/index.html"; \
	else \
		echo "⚠️  swag não encontrado. Instalando..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		$(HOME)/go/bin/swag init -g main.go -o docs --parseDependency --parseInternal; \
		echo "✅ Documentação Swagger gerada com sucesso!"; \
		echo "📖 Acesse: http://localhost:8080/swagger/index.html"; \
	fi