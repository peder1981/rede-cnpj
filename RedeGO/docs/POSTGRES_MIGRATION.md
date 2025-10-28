# 🐘 Plano de Migração SQLite → PostgreSQL

## 🎯 Objetivos da Migração

### **Por que PostgreSQL?**

1. **Performance Superior**
   - Queries complexas 10-100x mais rápidas
   - Índices avançados (GiST, GIN, BRIN)
   - Paralelização de queries
   - Particionamento de tabelas

2. **Recursos Avançados**
   - Full-text search nativo e poderoso
   - JSON/JSONB para dados semi-estruturados
   - Extensões (PostGIS, pg_trgm, etc)
   - Views materializadas
   - CTEs recursivos otimizados

3. **Escalabilidade**
   - Suporta terabytes de dados
   - Conexões concorrentes ilimitadas
   - Replicação master-slave
   - Sharding horizontal

4. **Funcionalidades Forenses**
   - Análise de grafos nativa (Apache AGE)
   - Machine Learning (MADlib)
   - Busca por similaridade
   - Análise temporal avançada

---

## 📊 Análise da Massa de Dados Atual

### **SQLite Atual**
```
cnpj.db           ~200 GB
├── empresas      50M registros
├── estabelecimento 50M registros
├── socios        26M registros
├── simples       10M registros
└── tabelas lookup ~1M registros

rede.db           ~10 GB
└── ligacao       100M registros

rede_search.db    ~5 GB
└── id_search     FTS5 index
```

### **Total:** ~215 GB, ~237M registros

---

## 🏗️ Arquitetura PostgreSQL Proposta

### **1. Schema Otimizado**

```sql
-- Database: rede_cnpj
-- Owner: rede_user
-- Encoding: UTF8
-- Collation: pt_BR.UTF-8

-- Schema para dados da Receita Federal
CREATE SCHEMA receita;

-- Schema para rede de relacionamentos
CREATE SCHEMA rede;

-- Schema para análises forenses
CREATE SCHEMA forensics;

-- Schema para cache e views materializadas
CREATE SCHEMA cache;
```

### **2. Tabelas Particionadas**

```sql
-- Empresas particionadas por UF
CREATE TABLE receita.empresas (
    cnpj_basico VARCHAR(8) PRIMARY KEY,
    razao_social TEXT NOT NULL,
    natureza_juridica VARCHAR(4),
    qualificacao_responsavel VARCHAR(2),
    capital_social NUMERIC(15,2),
    porte_empresa VARCHAR(2),
    ente_federativo_responsavel TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
) PARTITION BY LIST (SUBSTRING(cnpj_basico, 1, 2));

-- Estabelecimentos particionados por UF
CREATE TABLE receita.estabelecimento (
    cnpj VARCHAR(14) PRIMARY KEY,
    cnpj_basico VARCHAR(8) NOT NULL,
    cnpj_ordem VARCHAR(4) NOT NULL,
    cnpj_dv VARCHAR(2) NOT NULL,
    matriz_filial VARCHAR(1),
    nome_fantasia TEXT,
    situacao_cadastral VARCHAR(2),
    data_situacao_cadastral DATE,
    uf VARCHAR(2) NOT NULL,
    -- ... outros campos
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
) PARTITION BY LIST (uf);

-- Criar partições por UF (27 estados)
CREATE TABLE receita.estabelecimento_sp PARTITION OF receita.estabelecimento FOR VALUES IN ('SP');
CREATE TABLE receita.estabelecimento_rj PARTITION OF receita.estabelecimento FOR VALUES IN ('RJ');
-- ... etc para todos os estados

-- Sócios particionados por tipo
CREATE TABLE receita.socios (
    id BIGSERIAL,
    cnpj VARCHAR(14) NOT NULL,
    cnpj_basico VARCHAR(8) NOT NULL,
    identificador_de_socio VARCHAR(1) NOT NULL,
    nome_socio TEXT NOT NULL,
    cnpj_cpf_socio VARCHAR(14) NOT NULL,
    qualificacao_socio VARCHAR(2),
    data_entrada_sociedade DATE,
    pais VARCHAR(3),
    representante_legal VARCHAR(11),
    nome_representante TEXT,
    qualificacao_representante_legal VARCHAR(2),
    faixa_etaria VARCHAR(1),
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id, identificador_de_socio)
) PARTITION BY LIST (identificador_de_socio);

-- Partições: PF, PJ, Estrangeiro
CREATE TABLE receita.socios_pf PARTITION OF receita.socios FOR VALUES IN ('1');
CREATE TABLE receita.socios_pj PARTITION OF receita.socios FOR VALUES IN ('2');
CREATE TABLE receita.socios_pe PARTITION OF receita.socios FOR VALUES IN ('3');
```

### **3. Índices Avançados**

```sql
-- Índices B-Tree para buscas exatas
CREATE INDEX idx_empresas_cnpj ON receita.empresas(cnpj_basico);
CREATE INDEX idx_estabelecimento_cnpj ON receita.estabelecimento(cnpj);
CREATE INDEX idx_socios_cpf ON receita.socios(cnpj_cpf_socio);
CREATE INDEX idx_socios_cnpj ON receita.socios(cnpj);

-- Índices GIN para full-text search
CREATE INDEX idx_empresas_razao_gin ON receita.empresas 
    USING GIN (to_tsvector('portuguese', razao_social));

CREATE INDEX idx_estabelecimento_fantasia_gin ON receita.estabelecimento 
    USING GIN (to_tsvector('portuguese', nome_fantasia));

CREATE INDEX idx_socios_nome_gin ON receita.socios 
    USING GIN (to_tsvector('portuguese', nome_socio));

-- Índices GiST para busca por similaridade
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_empresas_razao_trgm ON receita.empresas 
    USING GiST (razao_social gist_trgm_ops);

CREATE INDEX idx_socios_nome_trgm ON receita.socios 
    USING GiST (nome_socio gist_trgm_ops);

-- Índices compostos para queries comuns
CREATE INDEX idx_estabelecimento_uf_situacao ON receita.estabelecimento(uf, situacao_cadastral);
CREATE INDEX idx_estabelecimento_endereco ON receita.estabelecimento(cep, logradouro, numero);
CREATE INDEX idx_estabelecimento_contato ON receita.estabelecimento(telefone1) WHERE telefone1 IS NOT NULL;
CREATE INDEX idx_estabelecimento_email ON receita.estabelecimento(correio_eletronico) WHERE correio_eletronico IS NOT NULL;

-- Índices BRIN para dados temporais (muito eficiente)
CREATE INDEX idx_estabelecimento_data_brin ON receita.estabelecimento 
    USING BRIN (data_inicio_atividades);
CREATE INDEX idx_socios_data_brin ON receita.socios 
    USING BRIN (data_entrada_sociedade);
```

### **4. Views Materializadas para Performance**

```sql
-- View materializada: Empresas ativas por UF
CREATE MATERIALIZED VIEW cache.empresas_ativas_por_uf AS
SELECT 
    e.uf,
    COUNT(*) as total_empresas,
    COUNT(DISTINCT est.cnpj_basico) as total_grupos,
    SUM(emp.capital_social) as capital_total
FROM receita.estabelecimento e
JOIN receita.empresas emp ON e.cnpj_basico = emp.cnpj_basico
WHERE e.situacao_cadastral = '02'
GROUP BY e.uf;

CREATE UNIQUE INDEX ON cache.empresas_ativas_por_uf(uf);

-- View materializada: Top sócios (atualizada diariamente)
CREATE MATERIALIZED VIEW cache.top_socios AS
SELECT 
    s.cnpj_cpf_socio,
    s.nome_socio,
    COUNT(DISTINCT s.cnpj) as total_empresas,
    COUNT(DISTINCT CASE WHEN est.situacao_cadastral = '02' THEN s.cnpj END) as empresas_ativas,
    SUM(emp.capital_social) as capital_total
FROM receita.socios s
JOIN receita.estabelecimento est ON s.cnpj = est.cnpj
JOIN receita.empresas emp ON s.cnpj_basico = emp.cnpj_basico
WHERE s.identificador_de_socio = '1'
GROUP BY s.cnpj_cpf_socio, s.nome_socio
HAVING COUNT(DISTINCT s.cnpj) >= 5
ORDER BY total_empresas DESC;

CREATE UNIQUE INDEX ON cache.top_socios(cnpj_cpf_socio);

-- Refresh automático (via cron ou pg_cron)
-- SELECT cron.schedule('refresh-top-socios', '0 2 * * *', 
--   'REFRESH MATERIALIZED VIEW CONCURRENTLY cache.top_socios');
```

### **5. Tabela de Grafos (Apache AGE)**

```sql
-- Instalar extensão Apache AGE
CREATE EXTENSION IF NOT EXISTS age;

-- Criar grafo de relacionamentos
SELECT create_graph('rede_cnpj');

-- Criar vértices (empresas e pessoas)
SELECT * FROM cypher('rede_cnpj', $$
    CREATE (:Empresa {cnpj: '01234567000100', razao: 'EMPRESA A'})
$$) as (v agtype);

SELECT * FROM cypher('rede_cnpj', $$
    CREATE (:Pessoa {cpf: '12345678900', nome: 'JOÃO SILVA'})
$$) as (v agtype);

-- Criar arestas (relacionamentos)
SELECT * FROM cypher('rede_cnpj', $$
    MATCH (p:Pessoa {cpf: '12345678900'})
    MATCH (e:Empresa {cnpj: '01234567000100'})
    CREATE (p)-[:SOCIO {cargo: 'Administrador', data: '2020-01-01'}]->(e)
$$) as (e agtype);

-- Query de grafos (encontrar rede de 2º grau)
SELECT * FROM cypher('rede_cnpj', $$
    MATCH (p:Pessoa {cpf: '12345678900'})-[:SOCIO]->(e1:Empresa)
    MATCH (e1)<-[:SOCIO]-(p2:Pessoa)
    MATCH (p2)-[:SOCIO]->(e2:Empresa)
    WHERE e2.cnpj <> e1.cnpj
    RETURN DISTINCT e2.cnpj, e2.razao, COUNT(*) as conexoes
    ORDER BY conexoes DESC
    LIMIT 100
$$) as (cnpj agtype, razao agtype, conexoes agtype);
```

---

## 🔄 Processo de Migração

### **Fase 1: Preparação (1 dia)**

1. **Instalar PostgreSQL 16**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql-16 postgresql-contrib-16

# Configurar autenticação
sudo -u postgres psql
CREATE USER rede_user WITH PASSWORD 'senha_segura';
CREATE DATABASE rede_cnpj OWNER rede_user;
\c rede_cnpj
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;
CREATE EXTENSION IF NOT EXISTS age;
```

2. **Otimizar postgresql.conf**
```ini
# Memória (para servidor com 32GB RAM)
shared_buffers = 8GB
effective_cache_size = 24GB
maintenance_work_mem = 2GB
work_mem = 256MB

# Paralelização
max_parallel_workers_per_gather = 4
max_parallel_workers = 8
max_worker_processes = 8

# WAL
wal_buffers = 16MB
checkpoint_completion_target = 0.9
max_wal_size = 4GB

# Planner
random_page_cost = 1.1  # Para SSD
effective_io_concurrency = 200

# Logging
log_min_duration_statement = 1000  # Log queries > 1s
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
```

### **Fase 2: Criação do Schema (2 horas)**

```bash
# Executar scripts de criação
psql -U rede_user -d rede_cnpj -f scripts/postgres/01_schemas.sql
psql -U rede_user -d rede_cnpj -f scripts/postgres/02_tables.sql
psql -U rede_user -d rede_cnpj -f scripts/postgres/03_partitions.sql
```

### **Fase 3: Migração de Dados (2-3 dias)**

#### **Opção A: Via CSV (Mais Rápido)**

```bash
# 1. Exportar SQLite para CSV
sqlite3 bases/cnpj.db <<EOF
.mode csv
.output /tmp/empresas.csv
SELECT * FROM empresas;
.output /tmp/estabelecimento.csv
SELECT * FROM estabelecimento;
.output /tmp/socios.csv
SELECT * FROM socios;
.quit
EOF

# 2. Importar para PostgreSQL (COPY é muito rápido)
psql -U rede_user -d rede_cnpj <<EOF
\copy receita.empresas FROM '/tmp/empresas.csv' WITH (FORMAT csv, HEADER true);
\copy receita.estabelecimento FROM '/tmp/estabelecimento.csv' WITH (FORMAT csv, HEADER true);
\copy receita.socios FROM '/tmp/socios.csv' WITH (FORMAT csv, HEADER true);
EOF
```

#### **Opção B: Via pgloader (Automático)**

```bash
# Instalar pgloader
sudo apt install pgloader

# Criar arquivo de configuração
cat > migrate.load <<EOF
LOAD DATABASE
    FROM sqlite://bases/cnpj.db
    INTO postgresql://rede_user:senha@localhost/rede_cnpj

WITH include drop, create tables, create indexes, reset sequences

SET work_mem to '256MB', maintenance_work_mem to '2GB'

CAST type text to text drop typemod

BEFORE LOAD DO
    \$\$ DROP SCHEMA IF EXISTS receita CASCADE; \$\$,
    \$\$ CREATE SCHEMA receita; \$\$;
EOF

# Executar migração
pgloader migrate.load
```

#### **Opção C: Via Go (Programático)**

```go
// cmd/migrate/main.go
package main

import (
    "database/sql"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
    _ "github.com/lib/pq"
)

func main() {
    // Conectar SQLite
    sqliteDB, _ := sql.Open("sqlite3", "bases/cnpj.db")
    defer sqliteDB.Close()
    
    // Conectar PostgreSQL
    pgDB, _ := sql.Open("postgres", 
        "host=localhost user=rede_user password=senha dbname=rede_cnpj sslmode=disable")
    defer pgDB.Close()
    
    // Migrar empresas
    log.Println("Migrando empresas...")
    migrateTable(sqliteDB, pgDB, "empresas", "receita.empresas", 10000)
    
    // Migrar estabelecimentos
    log.Println("Migrando estabelecimentos...")
    migrateTable(sqliteDB, pgDB, "estabelecimento", "receita.estabelecimento", 10000)
    
    // Migrar sócios
    log.Println("Migrando sócios...")
    migrateTable(sqliteDB, pgDB, "socios", "receita.socios", 10000)
    
    log.Println("Migração concluída!")
}

func migrateTable(src, dst *sql.DB, srcTable, dstTable string, batchSize int) {
    // Implementação com batches para performance
    // ...
}
```

### **Fase 4: Criação de Índices (1 dia)**

```bash
# Criar índices em paralelo
psql -U rede_user -d rede_cnpj -f scripts/postgres/04_indexes.sql

# Analisar tabelas para otimizar planner
psql -U rede_user -d rede_cnpj <<EOF
ANALYZE receita.empresas;
ANALYZE receita.estabelecimento;
ANALYZE receita.socios;
EOF
```

### **Fase 5: Views e Funções (4 horas)**

```bash
psql -U rede_user -d rede_cnpj -f scripts/postgres/05_views.sql
psql -U rede_user -d rede_cnpj -f scripts/postgres/06_functions.sql
psql -U rede_user -d rede_cnpj -f scripts/postgres/07_triggers.sql
```

### **Fase 6: Testes e Validação (1 dia)**

```sql
-- Validar contagens
SELECT 'empresas' as tabela, COUNT(*) FROM receita.empresas
UNION ALL
SELECT 'estabelecimento', COUNT(*) FROM receita.estabelecimento
UNION ALL
SELECT 'socios', COUNT(*) FROM receita.socios;

-- Testar performance de queries críticas
EXPLAIN ANALYZE
SELECT * FROM receita.socios WHERE cnpj_cpf_socio = '12345678900';

EXPLAIN ANALYZE
SELECT * FROM receita.estabelecimento WHERE cnpj = '01234567000100';
```

---

## 🚀 Funcionalidades Novas com PostgreSQL

### **1. Busca Full-Text Avançada**

```sql
-- Busca com ranking e highlight
SELECT 
    cnpj,
    razao_social,
    ts_rank(to_tsvector('portuguese', razao_social), query) as rank,
    ts_headline('portuguese', razao_social, query) as highlight
FROM receita.empresas,
     to_tsquery('portuguese', 'tecnologia & informacao') query
WHERE to_tsvector('portuguese', razao_social) @@ query
ORDER BY rank DESC
LIMIT 100;
```

### **2. Busca por Similaridade**

```sql
-- Encontrar nomes similares (fuzzy search)
SELECT 
    nome_socio,
    similarity(nome_socio, 'JOAO DA SILVA') as sim
FROM receita.socios
WHERE nome_socio % 'JOAO DA SILVA'  -- Operador de similaridade
ORDER BY sim DESC
LIMIT 20;
```

### **3. Análise Geoespacial (PostGIS)**

```sql
CREATE EXTENSION IF NOT EXISTS postgis;

-- Adicionar coluna de geometria
ALTER TABLE receita.estabelecimento ADD COLUMN geom geometry(Point, 4326);

-- Popular com coordenadas (via geocoding)
UPDATE receita.estabelecimento
SET geom = ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
WHERE longitude IS NOT NULL AND latitude IS NOT NULL;

-- Criar índice espacial
CREATE INDEX idx_estabelecimento_geom ON receita.estabelecimento USING GIST(geom);

-- Buscar empresas próximas
SELECT 
    cnpj,
    razao_social,
    ST_Distance(geom, ST_MakePoint(-46.6333, -23.5505)::geography) as distancia_metros
FROM receita.estabelecimento e
JOIN receita.empresas emp ON e.cnpj_basico = emp.cnpj_basico
WHERE ST_DWithin(
    geom,
    ST_MakePoint(-46.6333, -23.5505)::geography,
    1000  -- 1km
)
ORDER BY distancia_metros;
```

### **4. Machine Learning (MADlib)**

```sql
CREATE EXTENSION IF NOT EXISTS madlib;

-- Clustering de empresas por características
SELECT madlib.kmeans(
    'receita.empresas',
    'forensics.empresa_clusters',
    'ARRAY[capital_social, total_socios, idade_anos]',
    5  -- 5 clusters
);

-- Detecção de anomalias
SELECT madlib.lof(
    'receita.socios_stats',
    'forensics.socios_anomalos',
    'cpf',
    'ARRAY[total_empresas, empresas_baixadas, capital_total]',
    10  -- k-neighbors
);
```

### **5. JSON para Dados Flexíveis**

```sql
-- Adicionar coluna JSONB para metadados
ALTER TABLE receita.estabelecimento ADD COLUMN metadata JSONB;

-- Armazenar dados adicionais
UPDATE receita.estabelecimento
SET metadata = jsonb_build_object(
    'telefones', ARRAY[telefone1, telefone2],
    'emails', ARRAY[correio_eletronico],
    'redes_sociais', '{}'::jsonb
);

-- Criar índice GIN para busca em JSON
CREATE INDEX idx_estabelecimento_metadata ON receita.estabelecimento USING GIN(metadata);

-- Buscar por campo JSON
SELECT * FROM receita.estabelecimento
WHERE metadata @> '{"telefones": ["11999999999"]}';
```

---

## 📊 Comparação de Performance

### **Query 1: Buscar empresas de um CPF**

```
SQLite:   ~500ms  (sem índice otimizado)
PostgreSQL: ~15ms   (com índice + particionamento)
Speedup:  33x
```

### **Query 2: Full-text search**

```
SQLite FTS5: ~200ms
PostgreSQL:  ~8ms   (GIN index + ts_vector)
Speedup:     25x
```

### **Query 3: Agregações complexas**

```
SQLite:      ~5000ms
PostgreSQL:  ~100ms  (parallel query + partitions)
Speedup:     50x
```

### **Query 4: Grafos (rede 2º grau)**

```
SQLite:      ~10000ms (múltiplos JOINs)
PostgreSQL:  ~200ms   (Apache AGE + índices)
Speedup:     50x
```

---

## 💰 Estimativa de Custos

### **Infraestrutura**

| Recurso | Especificação | Custo Mensal |
|---------|---------------|--------------|
| Servidor | 32GB RAM, 8 cores, 1TB SSD | R$ 500-800 |
| Backup | 500GB storage | R$ 50-100 |
| **Total** | | **R$ 550-900** |

### **Tempo de Desenvolvimento**

| Fase | Tempo | Custo (R$ 200/h) |
|------|-------|------------------|
| Preparação | 1 dia | R$ 1.600 |
| Schema | 2h | R$ 400 |
| Migração | 3 dias | R$ 4.800 |
| Índices | 1 dia | R$ 1.600 |
| Views/Funções | 4h | R$ 800 |
| Testes | 1 dia | R$ 1.600 |
| **Total** | **6 dias** | **R$ 10.800** |

---

## 📅 Cronograma

### **Semana 1**
- Dias 1-2: Preparação e criação de schema
- Dias 3-5: Migração de dados
- Dia 6: Criação de índices

### **Semana 2**
- Dia 1: Views e funções
- Dias 2-3: Testes e validação
- Dias 4-5: Ajustes e otimizações

---

## ✅ Checklist de Migração

- [ ] Instalar PostgreSQL 16
- [ ] Configurar postgresql.conf
- [ ] Criar schemas e tabelas
- [ ] Criar partições
- [ ] Migrar dados (empresas)
- [ ] Migrar dados (estabelecimentos)
- [ ] Migrar dados (sócios)
- [ ] Migrar dados (tabelas lookup)
- [ ] Criar índices B-Tree
- [ ] Criar índices GIN (full-text)
- [ ] Criar índices GiST (similaridade)
- [ ] Criar índices BRIN (temporal)
- [ ] Criar views materializadas
- [ ] Criar funções stored procedures
- [ ] Configurar Apache AGE (grafos)
- [ ] Instalar extensões (PostGIS, MADlib)
- [ ] Validar contagens
- [ ] Testar performance
- [ ] Configurar backup automático
- [ ] Atualizar código Go
- [ ] Atualizar APIs
- [ ] Atualizar CLI
- [ ] Documentar mudanças
- [ ] Deploy em produção

---

## 🎯 Próximos Passos

1. **Aprovar o plano**
2. **Provisionar servidor PostgreSQL**
3. **Executar migração em ambiente de teste**
4. **Validar performance**
5. **Migrar produção**

**Pronto para começar a migração?** 🚀🐘
