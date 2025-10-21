#!/bin/bash

################################################################################
# Script de Inicialização do RedeCNPJ Go
# 
# Este script:
# - Verifica pré-requisitos
# - Instala dependências
# - Compila a aplicação
# - Inicia o servidor
# - Captura CTRL+C para encerrar graciosamente
#
# Uso: ./start.sh [opções]
#   -p, --port PORT       Porta do servidor (padrão: 5000)
#   -c, --config FILE     Arquivo de configuração (padrão: rede.ini)
#   -d, --dev             Modo desenvolvimento (recompila ao detectar mudanças)
#   -h, --help            Exibe esta ajuda
################################################################################

set -e  # Sair em caso de erro

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variáveis globais
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_NAME="rede-cnpj"
APP_PID=""
PORT=5000
CONFIG_FILE="rede.ini"
DEV_MODE=false

# Array para armazenar PIDs de processos iniciados
declare -a PIDS=()

################################################################################
# Funções de Utilidade
################################################################################

print_header() {
    echo -e "${BLUE}"
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║              RedeCNPJ Go - Script de Inicialização            ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERRO]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

################################################################################
# Função de Limpeza (executada ao sair)
################################################################################

cleanup() {
    echo ""
    log_info "Encerrando aplicação..."
    
    # Encerra o processo principal se estiver rodando
    if [ ! -z "$APP_PID" ] && kill -0 "$APP_PID" 2>/dev/null; then
        log_info "Encerrando servidor (PID: $APP_PID)..."
        kill -TERM "$APP_PID" 2>/dev/null || true
        
        # Aguarda até 5 segundos para encerramento gracioso
        for i in {1..5}; do
            if ! kill -0 "$APP_PID" 2>/dev/null; then
                break
            fi
            sleep 1
        done
        
        # Force kill se ainda estiver rodando
        if kill -0 "$APP_PID" 2>/dev/null; then
            log_warn "Forçando encerramento do servidor..."
            kill -9 "$APP_PID" 2>/dev/null || true
        fi
    fi
    
    # Encerra todos os processos filhos
    for pid in "${PIDS[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            log_info "Encerrando processo (PID: $pid)..."
            kill -TERM "$pid" 2>/dev/null || true
        fi
    done
    
    # Limpa processos zumbis
    wait 2>/dev/null || true
    
    log_info "Aplicação encerrada com sucesso!"
    exit 0
}

# Registra a função de limpeza para sinais de interrupção
trap cleanup SIGINT SIGTERM EXIT

################################################################################
# Verificações de Pré-requisitos
################################################################################

check_prerequisites() {
    log_step "Verificando pré-requisitos..."
    
    # Verifica se Go está instalado
    if ! command -v go &> /dev/null; then
        log_error "Go não está instalado!"
        log_info "Instale Go 1.21+ de: https://go.dev/dl/"
        exit 1
    fi
    
    # Verifica versão do Go
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go versão: $GO_VERSION"
    
    # Verifica se make está disponível (opcional)
    if command -v make &> /dev/null; then
        log_info "Make disponível: $(make --version | head -1)"
    else
        log_warn "Make não encontrado (opcional)"
    fi
    
    # Verifica se está no diretório correto
    if [ ! -f "go.mod" ]; then
        log_error "Arquivo go.mod não encontrado!"
        log_info "Execute este script do diretório RedeGO/"
        exit 1
    fi
    
    log_info "Pré-requisitos OK!"
}

################################################################################
# Criação de Diretórios
################################################################################

create_directories() {
    log_step "Criando diretórios necessários..."
    
    mkdir -p bases
    mkdir -p arquivos
    mkdir -p static
    mkdir -p templates
    mkdir -p logs
    
    log_info "Diretórios criados!"
}

################################################################################
# Verificação de Bancos de Dados
################################################################################

check_databases() {
    log_step "Verificando bancos de dados..."
    
    local missing_dbs=false
    
    # Verifica se existem bancos de dados de teste
    if [ ! -f "bases/cnpj_teste.db" ] && [ ! -f "bases/cnpj.db" ]; then
        log_warn "Nenhum banco de dados CNPJ encontrado em bases/"
        
        # Tenta copiar da versão Python
        if [ -f "../rede/bases/cnpj_teste.db" ]; then
            log_info "Copiando banco de teste da versão Python..."
            cp "../rede/bases/cnpj_teste.db" bases/ 2>/dev/null || true
            cp "../rede/bases/rede_teste.db" bases/ 2>/dev/null || true
        else
            missing_dbs=true
        fi
    fi
    
    if [ "$missing_dbs" = true ]; then
        log_warn "Bancos de dados não encontrados!"
        log_info "Copie os arquivos .db para o diretório bases/"
        log_info "Exemplo: cp ../rede/bases/*.db bases/"
    else
        log_info "Bancos de dados encontrados!"
    fi
}

################################################################################
# Instalação de Dependências
################################################################################

install_dependencies() {
    log_step "Instalando dependências Go..."
    
    if command -v make &> /dev/null && [ -f "Makefile" ]; then
        make deps
    else
        go mod download
        go mod tidy
    fi
    
    log_info "Dependências instaladas!"
}

################################################################################
# Compilação da Aplicação
################################################################################

build_application() {
    log_step "Compilando aplicação..."
    
    if command -v make &> /dev/null && [ -f "Makefile" ]; then
        make build
    else
        go build -o "$APP_NAME" ./cmd/server
    fi
    
    if [ ! -f "$APP_NAME" ]; then
        log_error "Falha na compilação!"
        exit 1
    fi
    
    chmod +x "$APP_NAME"
    log_info "Compilação concluída: $APP_NAME"
}

################################################################################
# Verificação de Configuração
################################################################################

check_configuration() {
    log_step "Verificando configuração..."
    
    if [ ! -f "$CONFIG_FILE" ]; then
        log_warn "Arquivo de configuração $CONFIG_FILE não encontrado!"
        log_info "Usando configuração padrão"
    else
        log_info "Configuração: $CONFIG_FILE"
    fi
}

################################################################################
# Inicialização do Servidor
################################################################################

start_server() {
    log_step "Iniciando servidor RedeCNPJ..."
    
    # Monta comando com opções
    CMD="./$APP_NAME"
    
    if [ -f "$CONFIG_FILE" ]; then
        CMD="$CMD -conf_file=$CONFIG_FILE"
    fi
    
    if [ "$PORT" != "5000" ]; then
        CMD="$CMD -porta_flask=$PORT"
    fi
    
    log_info "Comando: $CMD"
    log_info "Porta: $PORT"
    log_info ""
    log_info "Servidor iniciando..."
    log_info "Acesse: ${GREEN}http://127.0.0.1:$PORT/rede/${NC}"
    log_info ""
    log_warn "Pressione CTRL+C para encerrar"
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    # Inicia o servidor em background e captura o PID
    $CMD 2>&1 | tee logs/server.log &
    APP_PID=$!
    PIDS+=($APP_PID)
    
    # Aguarda um momento para verificar se iniciou corretamente
    sleep 2
    
    if ! kill -0 "$APP_PID" 2>/dev/null; then
        log_error "Falha ao iniciar o servidor!"
        log_info "Verifique os logs em: logs/server.log"
        exit 1
    fi
    
    log_info "Servidor rodando (PID: $APP_PID)"
    
    # Aguarda o processo terminar
    wait "$APP_PID"
}

################################################################################
# Modo Desenvolvimento (watch mode)
################################################################################

start_dev_mode() {
    log_step "Iniciando modo desenvolvimento..."
    log_info "Monitorando mudanças em arquivos .go"
    log_warn "Pressione CTRL+C para encerrar"
    echo ""
    
    # Função para recompilar e reiniciar
    restart_server() {
        if [ ! -z "$APP_PID" ] && kill -0 "$APP_PID" 2>/dev/null; then
            log_info "Reiniciando servidor..."
            kill -TERM "$APP_PID" 2>/dev/null || true
            wait "$APP_PID" 2>/dev/null || true
        fi
        
        build_application
        start_server &
    }
    
    # Inicia servidor pela primeira vez
    build_application
    start_server &
    
    # Monitora mudanças (requer inotify-tools)
    if command -v inotifywait &> /dev/null; then
        while true; do
            inotifywait -r -e modify,create,delete --include '\.go$' . 2>/dev/null
            log_info "Mudança detectada, recompilando..."
            restart_server
        done
    else
        log_warn "inotifywait não encontrado, modo watch desabilitado"
        log_info "Instale: sudo apt-get install inotify-tools"
        wait "$APP_PID"
    fi
}

################################################################################
# Exibição de Ajuda
################################################################################

show_help() {
    cat << EOF
Uso: $0 [opções]

Opções:
  -p, --port PORT       Porta do servidor (padrão: 5000)
  -c, --config FILE     Arquivo de configuração (padrão: rede.ini)
  -d, --dev             Modo desenvolvimento (recompila ao detectar mudanças)
  -h, --help            Exibe esta ajuda

Exemplos:
  $0                    # Inicia com configurações padrão
  $0 -p 8080            # Inicia na porta 8080
  $0 -d                 # Inicia em modo desenvolvimento
  $0 -c custom.ini      # Usa arquivo de configuração customizado

Atalhos:
  CTRL+C                Encerra o servidor graciosamente

Logs:
  logs/server.log       Log do servidor

EOF
}

################################################################################
# Processamento de Argumentos
################################################################################

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -p|--port)
                PORT="$2"
                shift 2
                ;;
            -c|--config)
                CONFIG_FILE="$2"
                shift 2
                ;;
            -d|--dev)
                DEV_MODE=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "Opção desconhecida: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

################################################################################
# Função Principal
################################################################################

main() {
    # Muda para o diretório do script
    cd "$SCRIPT_DIR"
    
    # Exibe cabeçalho
    print_header
    
    # Processa argumentos
    parse_arguments "$@"
    
    # Executa verificações e preparações
    check_prerequisites
    create_directories
    check_databases
    install_dependencies
    build_application
    check_configuration
    
    echo ""
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    log_info "Todas as verificações concluídas com sucesso!"
    log_info "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    # Inicia servidor (modo normal ou desenvolvimento)
    if [ "$DEV_MODE" = true ]; then
        start_dev_mode
    else
        start_server
    fi
}

################################################################################
# Execução
################################################################################

main "$@"
