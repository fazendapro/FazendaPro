# Makefile para FazendaPro API

.PHONY: help test test-coverage test-unit test-handlers clean install-deps run build migrate-docker

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
	go test -v -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./$(TEST_DIR)/...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	go tool cover -func=$(COVERAGE_DIR)/coverage.out
	@echo "📁 Relatório HTML gerado: $(COVERAGE_DIR)/coverage.html"

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