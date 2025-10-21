#!/bin/bash

################################################################################
# Script de Reinicialização do RedeCNPJ Go
# 
# Este script encerra e reinicia a aplicação
#
# Uso: ./restart.sh [opções do start.sh]
################################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║         RedeCNPJ Go - Reiniciando Aplicação                   ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Encerra processos existentes
if [ -f "$SCRIPT_DIR/stop.sh" ]; then
    bash "$SCRIPT_DIR/stop.sh"
else
    echo "Encerrando processos existentes..."
    pkill -f "rede-cnpj" || true
    sleep 2
fi

echo ""
echo "Aguardando 2 segundos..."
sleep 2
echo ""

# Inicia novamente
if [ -f "$SCRIPT_DIR/start.sh" ]; then
    bash "$SCRIPT_DIR/start.sh" "$@"
else
    echo "Erro: start.sh não encontrado!"
    exit 1
fi
