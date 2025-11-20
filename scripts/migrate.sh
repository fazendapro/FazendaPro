#!/bin/bash

# Script para executar migrações
# Uso: ./scripts/migrate.sh [up|down|status]

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para imprimir mensagens coloridas
print_message() {
    echo -e "${GREEN}[MIGRATE]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

if [ $# -eq 0 ]; then
    print_error "Uso: $0 [up|down|status]"
    echo "  up     - Executar migrações pendentes"
    echo "  down   - Reverter última migração"
    echo "  status - Verificar status das migrações"
    exit 1
fi

COMMAND=$1

case $COMMAND in
    "up")
        print_message "Executando migrações..."
        go run main.go migrate
        print_message "Migrações executadas com sucesso!"
        ;;
    "down")
        print_warning "Revertendo última migração..."
        print_warning "Rollback não implementado ainda"
        ;;
    "status")
        print_message "Verificando status das migrações..."
        print_warning "Status não implementado ainda"
        ;;
    *)
        print_error "Comando inválido: $COMMAND"
        echo "Comandos disponíveis: up, down, status"
        exit 1
        ;;
esac 