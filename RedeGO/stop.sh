#!/bin/bash

################################################################################
# Script de Parada do RedeCNPJ Go
# 
# Este script encerra todos os processos do RedeCNPJ em execução
#
# Uso: ./stop.sh
################################################################################

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERRO]${NC} $1"
}

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║           RedeCNPJ Go - Encerrando Aplicação                  ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Procura processos do rede-cnpj
PIDS=$(pgrep -f "rede-cnpj" || true)

if [ -z "$PIDS" ]; then
    log_info "Nenhum processo rede-cnpj em execução"
    exit 0
fi

log_info "Processos encontrados: $PIDS"

# Encerra graciosamente
for pid in $PIDS; do
    if kill -0 "$pid" 2>/dev/null; then
        log_info "Encerrando processo $pid..."
        kill -TERM "$pid" 2>/dev/null || true
    fi
done

# Aguarda até 5 segundos
sleep 2

# Verifica se ainda há processos rodando
REMAINING=$(pgrep -f "rede-cnpj" || true)

if [ ! -z "$REMAINING" ]; then
    log_warn "Forçando encerramento dos processos restantes..."
    for pid in $REMAINING; do
        kill -9 "$pid" 2>/dev/null || true
    done
fi

log_info "Todos os processos foram encerrados!"
