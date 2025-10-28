#!/bin/bash

################################################################################
# RedeCNPJ - PostgreSQL Setup Script
################################################################################
# Este script automatiza a criação completa do database PostgreSQL
################################################################################

set -e  # Sair em caso de erro

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variáveis
POSTGRES_USER="postgres"
DB_NAME="rede_cnpj"
DB_USER="rede_user"
DB_PASSWORD="rede_cnpj_2025"
DB_PORT="5433"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

################################################################################
# Funções
################################################################################

print_header() {
    echo -e "${BLUE}"
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║                                                                ║"
    echo "║     🐘 RedeCNPJ - PostgreSQL Setup                            ║"
    echo "║                                                                ║"
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
# Verificações
################################################################################

check_postgres() {
    log_step "Verificando PostgreSQL..."
    
    if ! command -v psql &> /dev/null; then
        log_error "PostgreSQL não está instalado!"
        log_info "Instale com: sudo apt install postgresql postgresql-contrib"
        exit 1
    fi
    
    PG_VERSION=$(psql --version | awk '{print $3}' | cut -d. -f1)
    log_info "PostgreSQL versão: $PG_VERSION"
    
    if [ "$PG_VERSION" -lt 12 ]; then
        log_warn "Recomendado PostgreSQL 12 ou superior"
    fi
}

check_postgres_running() {
    log_step "Verificando se PostgreSQL está rodando..."
    
    if ! sudo systemctl is-active --quiet postgresql; then
        log_warn "PostgreSQL não está rodando. Iniciando..."
        sudo systemctl start postgresql
        sleep 2
    fi
    
    log_info "PostgreSQL está rodando"
}

################################################################################
# Criação do Database
################################################################################

create_database() {
    log_step "Criando database e usuário..."
    
    # Criar usuário e database
    sudo -u postgres psql <<EOF
-- Criar usuário se não existir
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '$DB_USER') THEN
        CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';
        RAISE NOTICE 'Usuário $DB_USER criado';
    ELSE
        RAISE NOTICE 'Usuário $DB_USER já existe';
    END IF;
END
\$\$;

-- Dar privilégios
ALTER USER $DB_USER WITH CREATEDB CREATEROLE;

-- Criar database se não existir
SELECT 'Database já existe' WHERE EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME')
UNION ALL
SELECT 'Criando database...' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME');

-- Tentar criar (pode falhar se já existir)
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME') THEN
        PERFORM dblink_exec('dbname=' || current_database(), 
            'CREATE DATABASE $DB_NAME OWNER $DB_USER ENCODING UTF8');
    END IF;
EXCEPTION WHEN OTHERS THEN
    -- Database já existe, ignorar erro
    NULL;
END
\$\$;
EOF

    # Alternativa: criar database diretamente
    sudo -u postgres createdb -O $DB_USER $DB_NAME 2>/dev/null || log_info "Database já existe"
    
    log_info "Database $DB_NAME criado/verificado"
}

################################################################################
# Executar Scripts SQL
################################################################################

execute_sql_scripts() {
    log_step "Executando scripts SQL..."
    
    # 1. Schemas
    log_info "Criando schemas..."
    PGPASSWORD=$DB_PASSWORD psql -h localhost -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$SCRIPT_DIR/01_schemas.sql"
    
    # 2. Tabelas
    log_info "Criando tabelas..."
    PGPASSWORD=$DB_PASSWORD psql -h localhost -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$SCRIPT_DIR/02_tables.sql"
    
    # 3. Índices (se existir)
    if [ -f "$SCRIPT_DIR/03_indexes.sql" ]; then
        log_info "Criando índices..."
        PGPASSWORD=$DB_PASSWORD psql -U $DB_USER -d $DB_NAME -f "$SCRIPT_DIR/03_indexes.sql"
    fi
    
    # 4. Views (se existir)
    if [ -f "$SCRIPT_DIR/04_views.sql" ]; then
        log_info "Criando views..."
        PGPASSWORD=$DB_PASSWORD psql -U $DB_USER -d $DB_NAME -f "$SCRIPT_DIR/04_views.sql"
    fi
    
    log_info "Scripts executados com sucesso!"
}

################################################################################
# Verificação Final
################################################################################

verify_setup() {
    log_step "Verificando instalação..."
    
    PGPASSWORD=$DB_PASSWORD psql -h localhost -p $DB_PORT -U $DB_USER -d $DB_NAME <<EOF
-- Listar schemas
SELECT schema_name FROM information_schema.schemata 
WHERE schema_name IN ('receita', 'rede', 'forensics', 'cache');

-- Listar tabelas
SELECT schemaname, tablename 
FROM pg_tables 
WHERE schemaname IN ('receita', 'rede', 'forensics', 'cache')
ORDER BY schemaname, tablename;

-- Listar extensões
SELECT extname, extversion FROM pg_extension;
EOF
    
    log_info "Verificação concluída!"
}

################################################################################
# Criar arquivo .pgpass para não pedir senha
################################################################################

create_pgpass() {
    log_step "Criando arquivo .pgpass..."
    
    PGPASS_FILE="$HOME/.pgpass"
    PGPASS_LINE="localhost:$DB_PORT:$DB_NAME:$DB_USER:$DB_PASSWORD"
    
    # Criar ou atualizar .pgpass
    if [ -f "$PGPASS_FILE" ]; then
        # Remover linha antiga se existir
        grep -v "$DB_NAME:$DB_USER" "$PGPASS_FILE" > "$PGPASS_FILE.tmp" || true
        echo "$PGPASS_LINE" >> "$PGPASS_FILE.tmp"
        mv "$PGPASS_FILE.tmp" "$PGPASS_FILE"
    else
        echo "$PGPASS_LINE" > "$PGPASS_FILE"
    fi
    
    chmod 600 "$PGPASS_FILE"
    log_info "Arquivo .pgpass configurado"
}

################################################################################
# Informações de Conexão
################################################################################

show_connection_info() {
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                                                                ║${NC}"
    echo -e "${GREEN}║     ✅ PostgreSQL configurado com sucesso!                     ║${NC}"
    echo -e "${GREEN}║                                                                ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "📋 Informações de Conexão:"
    echo "   Host:     localhost"
    echo "   Porta:    $DB_PORT"
    echo "   Database: $DB_NAME"
    echo "   Usuário:  $DB_USER"
    echo "   Senha:    $DB_PASSWORD"
    echo ""
    echo "🔌 Conectar via psql:"
    echo "   psql -h localhost -p $DB_PORT -U $DB_USER -d $DB_NAME"
    echo ""
    echo "🔌 String de conexão:"
    echo "   postgresql://$DB_USER:$DB_PASSWORD@localhost:$DB_PORT/$DB_NAME"
    echo ""
    echo "📊 Próximos passos:"
    echo "   1. Migrar dados do SQLite"
    echo "   2. Criar índices adicionais"
    echo "   3. Atualizar código Go para usar PostgreSQL"
    echo ""
}

################################################################################
# Main
################################################################################

main() {
    print_header
    
    check_postgres
    check_postgres_running
    create_database
    execute_sql_scripts
    create_pgpass
    verify_setup
    show_connection_info
}

# Executar
main "$@"
