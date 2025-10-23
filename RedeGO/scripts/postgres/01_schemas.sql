-- ============================================================================
-- RedeCNPJ - PostgreSQL Schema Creation
-- ============================================================================
-- Cria os schemas principais do banco de dados
-- ============================================================================

-- Drop schemas se existirem (cuidado em produção!)
DROP SCHEMA IF EXISTS receita CASCADE;
DROP SCHEMA IF EXISTS rede CASCADE;
DROP SCHEMA IF EXISTS forensics CASCADE;
DROP SCHEMA IF EXISTS cache CASCADE;

-- Criar schemas
CREATE SCHEMA receita;
COMMENT ON SCHEMA receita IS 'Dados da Receita Federal (empresas, sócios, estabelecimentos)';

CREATE SCHEMA rede;
COMMENT ON SCHEMA rede IS 'Rede de relacionamentos e ligações';

CREATE SCHEMA forensics;
COMMENT ON SCHEMA forensics IS 'Análises forenses, scores de risco, alertas';

CREATE SCHEMA cache;
COMMENT ON SCHEMA cache IS 'Views materializadas e cache de queries';

-- Instalar extensões necessárias
CREATE EXTENSION IF NOT EXISTS pg_trgm;
COMMENT ON EXTENSION pg_trgm IS 'Busca por similaridade de texto';

CREATE EXTENSION IF NOT EXISTS btree_gin;
COMMENT ON EXTENSION btree_gin IS 'Índices GIN para tipos B-tree';

CREATE EXTENSION IF NOT EXISTS unaccent;
COMMENT ON EXTENSION unaccent IS 'Remove acentos para busca';

-- Extensões opcionais (comentar se não disponíveis)
-- CREATE EXTENSION IF NOT EXISTS postgis;
-- COMMENT ON EXTENSION postgis IS 'Análise geoespacial';

-- CREATE EXTENSION IF NOT EXISTS age;
-- COMMENT ON EXTENSION age IS 'Apache AGE - Graph database';

-- CREATE EXTENSION IF NOT EXISTS madlib;
-- COMMENT ON EXTENSION madlib IS 'Machine Learning';

-- Configurar search_path padrão
ALTER DATABASE rede_cnpj SET search_path TO receita, rede, forensics, cache, public;

-- Criar função para atualizar timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_updated_at_column() IS 'Atualiza automaticamente coluna updated_at';

-- Mensagem de sucesso
DO $$
BEGIN
    RAISE NOTICE 'Schemas criados com sucesso!';
    RAISE NOTICE 'Extensões instaladas: pg_trgm, btree_gin, unaccent';
END $$;
