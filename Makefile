# Makefile para FazendaPro API

.PHONY: help test test-coverage test-unit test-handlers clean install-deps run build

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
	@echo ""
	@echo "🔧 Desenvolvimento:"
	@echo "  run             - Executar a aplicação"
	@echo "  build           - Compilar a aplicação"
	@echo "  clean           - Limpar arquivos temporários"
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