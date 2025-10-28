#!/bin/bash

################################################################################
# Script para Monitorar Migração PostgreSQL
################################################################################

POSTGRES_URL="postgresql://rede_user:rede_cnpj_2025@localhost:5433/rede_cnpj?sslmode=disable"
LOG_FILE="../../migrate.log"

# Cores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

show_status() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                                                                ║${NC}"
    echo -e "${BLUE}║     📊 Status da Migração PostgreSQL                          ║${NC}"
    echo -e "${BLUE}║                                                                ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    
    # Verificar se processo está rodando
    if pgrep -f "rede-cnpj-migrate" > /dev/null; then
        echo -e "${GREEN}✅ Processo de migração: RODANDO${NC}"
        PID=$(pgrep -f "rede-cnpj-migrate")
        echo -e "   PID: $PID"
    else
        echo -e "${RED}❌ Processo de migração: PARADO${NC}"
    fi
    
    echo ""
    echo -e "${YELLOW}📈 Contagem de Registros:${NC}"
    echo ""
    
    # Conectar ao PostgreSQL e contar registros
    PGPASSWORD=rede_cnpj_2025 psql -h localhost -p 5433 -U rede_user -d rede_cnpj -t -c "
        SELECT 
            '   Empresas:         ' || LPAD(COUNT(*)::text, 12) as count
        FROM receita.empresas
        UNION ALL
        SELECT 
            '   Estabelecimentos: ' || LPAD(COUNT(*)::text, 12)
        FROM receita.estabelecimento
        UNION ALL
        SELECT 
            '   Sócios:           ' || LPAD(COUNT(*)::text, 12)
        FROM receita.socios
        UNION ALL
        SELECT 
            '   Simples:          ' || LPAD(COUNT(*)::text, 12)
        FROM receita.simples;
    " 2>/dev/null || echo -e "${RED}   Erro ao conectar ao PostgreSQL${NC}"
    
    echo ""
}

show_logs() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                                                                ║${NC}"
    echo -e "${BLUE}║     📋 Últimas 20 Linhas do Log                               ║${NC}"
    echo -e "${BLUE}║                                                                ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    
    if [ -f "$LOG_FILE" ]; then
        tail -20 "$LOG_FILE"
    else
        echo -e "${RED}❌ Arquivo de log não encontrado: $LOG_FILE${NC}"
    fi
}

follow_logs() {
    echo -e "${BLUE}📋 Acompanhando logs em tempo real (Ctrl+C para sair)...${NC}"
    echo ""
    
    if [ -f "$LOG_FILE" ]; then
        tail -f "$LOG_FILE"
    else
        echo -e "${RED}❌ Arquivo de log não encontrado: $LOG_FILE${NC}"
    fi
}

start_migration() {
    if pgrep -f "rede-cnpj-migrate" > /dev/null; then
        echo -e "${YELLOW}⚠️  Migração já está rodando!${NC}"
        return
    fi
    
    echo -e "${GREEN}🚀 Iniciando migração em background...${NC}"
    cd ../..
    nohup ./bin/rede-cnpj-migrate > migrate.log 2>&1 &
    PID=$!
    echo -e "${GREEN}✅ Migração iniciada com PID: $PID${NC}"
    echo -e "   Log: migrate.log"
}

stop_migration() {
    if ! pgrep -f "rede-cnpj-migrate" > /dev/null; then
        echo -e "${YELLOW}⚠️  Migração não está rodando${NC}"
        return
    fi
    
    echo -e "${RED}🛑 Parando migração...${NC}"
    pkill -f "rede-cnpj-migrate"
    sleep 2
    
    if pgrep -f "rede-cnpj-migrate" > /dev/null; then
        echo -e "${RED}❌ Processo não parou. Forçando...${NC}"
        pkill -9 -f "rede-cnpj-migrate"
    fi
    
    echo -e "${GREEN}✅ Migração parada${NC}"
}

show_help() {
    echo "Uso: $0 [comando]"
    echo ""
    echo "Comandos:"
    echo "  status    - Mostra status da migração e contagem de registros"
    echo "  logs      - Mostra últimas 20 linhas do log"
    echo "  follow    - Acompanha logs em tempo real"
    echo "  start     - Inicia migração em background"
    echo "  stop      - Para a migração"
    echo "  restart   - Reinicia a migração"
    echo "  help      - Mostra esta ajuda"
    echo ""
}

# Main
case "${1:-status}" in
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
    follow)
        follow_logs
        ;;
    start)
        start_migration
        ;;
    stop)
        stop_migration
        ;;
    restart)
        stop_migration
        sleep 2
        start_migration
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo -e "${RED}❌ Comando inválido: $1${NC}"
        echo ""
        show_help
        exit 1
        ;;
esac
