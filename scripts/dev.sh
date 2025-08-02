#!/bin/bash

# Script para desenvolvimento local
# Uso: ./scripts/dev.sh [start|stop|reset|logs]

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_message() {
    echo -e "${GREEN}[DEV]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

# Verificar se o comando foi fornecido
if [ $# -eq 0 ]; then
    print_error "Uso: $0 [start|stop|reset|logs|db]"
    echo "  start - Inicia ambiente de desenvolvimento"
    echo "  stop  - Para o ambiente"
    echo "  reset - Reseta o banco de dados"
    echo "  logs  - Mostra logs"
    echo "  db    - Conecta ao banco de dados"
    exit 1
fi

COMMAND=$1

case $COMMAND in
    "start")
        print_message "Iniciando ambiente de desenvolvimento..."
        
        make dev-setup
        
        print_info "Iniciando PostgreSQL..."
        docker-compose up -d postgres
        
        print_info "Aguardando banco de dados..."
        sleep 10
        
        print_info "Executando migrações..."
        make migrate-docker
        
        print_info "Iniciando aplicação..."
        print_message "Aplicação disponível em: http://localhost:8080"
        print_message "Banco de dados em: localhost:5432"
        
        go run main.go
        ;;
    "stop")
        print_message "Parando ambiente..."
        docker-compose down
        print_message "Ambiente parado!"
        ;;
    "reset")
        print_warning "Resetando banco de dados..."
        docker-compose down -v
        docker-compose up -d postgres
        sleep 10
        make migrate-docker
        print_message "Banco resetado!"
        ;;
    "logs")
        print_info "Mostrando logs..."
        docker-compose logs -f
        ;;
    "db")
        print_info "Conectando ao banco de dados..."
        docker-compose exec postgres psql -U fazendapro_user -d fazendapro
        ;;
    *)
        print_error "Comando inválido: $COMMAND"
        echo "Comandos disponíveis: start, stop, reset, logs, db"
        exit 1
        ;;
esac 