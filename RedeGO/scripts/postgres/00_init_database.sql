-- ============================================================================
-- RedeCNPJ - PostgreSQL Database Initialization
-- ============================================================================
-- Este script deve ser executado como superusuário (postgres)
-- Cria o database, usuário e configurações iniciais
-- ============================================================================

-- Conectar como postgres:
-- psql -U postgres

-- ============================================================================
-- 1. CRIAR USUÁRIO
-- ============================================================================

-- Verificar se usuário existe
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'rede_user') THEN
        CREATE USER rede_user WITH PASSWORD 'rede_cnpj_2025';
        RAISE NOTICE 'Usuário rede_user criado';
    ELSE
        RAISE NOTICE 'Usuário rede_user já existe';
    END IF;
END
$$;

-- Dar privilégios ao usuário
ALTER USER rede_user WITH CREATEDB;
ALTER USER rede_user WITH CREATEROLE;

COMMENT ON ROLE rede_user IS 'Usuário principal do sistema RedeCNPJ';

-- ============================================================================
-- 2. CRIAR DATABASE
-- ============================================================================

-- Verificar se database existe
SELECT 'Database rede_cnpj já existe' AS status
WHERE EXISTS (SELECT FROM pg_database WHERE datname = 'rede_cnpj')
UNION ALL
SELECT 'Criando database rede_cnpj...' AS status
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'rede_cnpj');

-- Criar database (executar fora de transação)
-- Nota: Este comando pode falhar se o database já existir
CREATE DATABASE rede_cnpj
    WITH 
    OWNER = rede_user
    ENCODING = 'UTF8'
    LC_COLLATE = 'pt_BR.UTF-8'
    LC_CTYPE = 'pt_BR.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    TEMPLATE = template0;

COMMENT ON DATABASE rede_cnpj IS 'Base de dados RedeCNPJ - Dados da Receita Federal';

-- ============================================================================
-- 3. CONECTAR AO DATABASE E CONFIGURAR
-- ============================================================================

-- Agora conectar ao database criado:
-- \c rede_cnpj

-- Dar todos os privilégios ao usuário
GRANT ALL PRIVILEGES ON DATABASE rede_cnpj TO rede_user;

-- Configurar search_path padrão
ALTER DATABASE rede_cnpj SET search_path TO receita, rede, forensics, cache, public;

-- Configurações de performance
ALTER DATABASE rede_cnpj SET shared_buffers TO '8GB';
ALTER DATABASE rede_cnpj SET effective_cache_size TO '24GB';
ALTER DATABASE rede_cnpj SET maintenance_work_mem TO '2GB';
ALTER DATABASE rede_cnpj SET work_mem TO '256MB';
ALTER DATABASE rede_cnpj SET random_page_cost TO '1.1';

-- Configurações de logging
ALTER DATABASE rede_cnpj SET log_min_duration_statement TO '1000';
ALTER DATABASE rede_cnpj SET log_statement TO 'mod';

-- ============================================================================
-- 4. MENSAGEM FINAL
-- ============================================================================

DO $$
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '╔════════════════════════════════════════════════════════════════╗';
    RAISE NOTICE '║                                                                ║';
    RAISE NOTICE '║     ✅ Database rede_cnpj criado com sucesso!                  ║';
    RAISE NOTICE '║                                                                ║';
    RAISE NOTICE '╚════════════════════════════════════════════════════════════════╝';
    RAISE NOTICE '';
    RAISE NOTICE 'Próximos passos:';
    RAISE NOTICE '1. Conectar ao database: \c rede_cnpj';
    RAISE NOTICE '2. Executar: psql -U rede_user -d rede_cnpj -f 01_schemas.sql';
    RAISE NOTICE '3. Executar: psql -U rede_user -d rede_cnpj -f 02_tables.sql';
    RAISE NOTICE '';
    RAISE NOTICE 'Credenciais:';
    RAISE NOTICE '  Usuário: rede_user';
    RAISE NOTICE '  Senha: rede_cnpj_2025';
    RAISE NOTICE '  Database: rede_cnpj';
    RAISE NOTICE '';
END $$;
