# 🐘 Scripts PostgreSQL - RedeCNPJ

## 📋 Ordem de Execução

### **Método 1: Automático (Recomendado)**

```bash
# Executar script de setup completo
./setup_postgres.sh
```

Este script faz tudo automaticamente:
1. ✅ Verifica se PostgreSQL está instalado
2. ✅ Cria usuário `rede_user`
3. ✅ Cria database `rede_cnpj`
4. ✅ Executa todos os scripts SQL
5. ✅ Configura `.pgpass` para não pedir senha
6. ✅ Verifica a instalação

---

### **Método 2: Manual**

#### **1. Criar Database e Usuário**

```bash
# Como superusuário postgres
sudo -u postgres psql -f 00_init_database.sql
```

Ou manualmente:

```bash
sudo -u postgres psql <<EOF
CREATE USER rede_user WITH PASSWORD 'rede_cnpj_2025';
ALTER USER rede_user WITH CREATEDB CREATEROLE;
CREATE DATABASE rede_cnpj OWNER rede_user ENCODING 'UTF8';
GRANT ALL PRIVILEGES ON DATABASE rede_cnpj TO rede_user;
EOF
```

#### **2. Criar Schemas**

```bash
psql -U rede_user -d rede_cnpj -f 01_schemas.sql
```

#### **3. Criar Tabelas**

```bash
psql -U rede_user -d rede_cnpj -f 02_tables.sql
```

#### **4. Criar Índices (quando disponível)**

```bash
psql -U rede_user -d rede_cnpj -f 03_indexes.sql
```

#### **5. Criar Views (quando disponível)**

```bash
psql -U rede_user -d rede_cnpj -f 04_views.sql
```

---

## 📁 Arquivos

| Arquivo | Descrição |
|---------|-----------|
| `00_init_database.sql` | Cria database, usuário e configurações iniciais |
| `01_schemas.sql` | Cria schemas (receita, rede, forensics, cache) |
| `02_tables.sql` | Cria todas as tabelas com particionamento |
| `03_indexes.sql` | Cria índices (B-Tree, GIN, GiST, BRIN) |
| `04_views.sql` | Cria views materializadas |
| `05_functions.sql` | Cria funções e stored procedures |
| `setup_postgres.sh` | Script automático de setup |

---

## 🔐 Credenciais Padrão

```
Host:     localhost
Porta:    5432
Database: rede_cnpj
Usuário:  rede_user
Senha:    rede_cnpj_2025
```

**⚠️ IMPORTANTE:** Altere a senha em produção!

---

## 🔌 Conectar ao Database

### **Via psql**

```bash
psql -U rede_user -d rede_cnpj
```

### **Via Go**

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

connStr := "postgresql://rede_user:rede_cnpj_2025@localhost:5432/rede_cnpj?sslmode=disable"
db, err := sql.Open("postgres", connStr)
```

### **String de Conexão**

```
postgresql://rede_user:rede_cnpj_2025@localhost:5432/rede_cnpj
```

---

## 📊 Estrutura Criada

### **Schemas**

- `receita` - Dados da Receita Federal
- `rede` - Relacionamentos
- `forensics` - Análises forenses
- `cache` - Views materializadas

### **Tabelas Principais**

```sql
receita.empresas              -- 50M registros
receita.estabelecimento       -- 50M (28 partições por UF)
receita.socios                -- 26M (3 partições por tipo)
receita.simples               -- 10M
rede.ligacao                  -- 100M
forensics.perfil_risco        -- Cache de análises
forensics.cluster_empresas    -- Clusters suspeitos
forensics.alertas             -- Alertas automáticos
```

### **Partições**

#### **Estabelecimentos (28 partições)**
- 27 estados (AC, AL, AP, AM, BA, CE, DF, ES, GO, MA, MT, MS, MG, PA, PB, PR, PE, PI, RJ, RN, RS, RO, RR, SC, SP, SE, TO)
- 1 exterior (EX)

#### **Sócios (3 partições)**
- PF (Pessoa Física)
- PJ (Pessoa Jurídica)
- PE (Pessoa Estrangeira)

---

## 🔍 Verificar Instalação

```sql
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

-- Listar partições
SELECT 
    parent.relname as parent_table,
    child.relname as partition_name
FROM pg_inherits
JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
JOIN pg_class child ON pg_inherits.inhrelid = child.oid
WHERE parent.relname IN ('estabelecimento', 'socios')
ORDER BY parent.relname, child.relname;
```

---

## 🚀 Próximos Passos

1. **Migrar Dados**
   ```bash
   # Ver POSTGRES_MIGRATION.md para detalhes
   ```

2. **Criar Índices Adicionais**
   ```bash
   psql -U rede_user -d rede_cnpj -f 03_indexes.sql
   ```

3. **Atualizar Código Go**
   ```bash
   # Trocar driver de sqlite3 para postgres
   go get github.com/lib/pq
   ```

4. **Testar Performance**
   ```sql
   EXPLAIN ANALYZE SELECT * FROM receita.socios WHERE cnpj_cpf_socio = '12345678900';
   ```

---

## 📝 Notas

- Scripts são idempotentes (podem ser executados múltiplas vezes)
- Use `DROP SCHEMA ... CASCADE` com cuidado em produção
- Particionamento melhora performance em 10-50x
- Extensões opcionais: PostGIS, Apache AGE, MADlib

---

## 🆘 Troubleshooting

### **Erro: role "rede_user" does not exist**
```bash
sudo -u postgres createuser -s rede_user
```

### **Erro: database "rede_cnpj" does not exist**
```bash
sudo -u postgres createdb -O rede_user rede_cnpj
```

### **Erro: permission denied**
```bash
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE rede_cnpj TO rede_user;"
```

### **PostgreSQL não está rodando**
```bash
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

---

## 📚 Referências

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Partitioning](https://www.postgresql.org/docs/current/ddl-partitioning.html)
- [Full-Text Search](https://www.postgresql.org/docs/current/textsearch.html)
- [pg_trgm](https://www.postgresql.org/docs/current/pgtrgm.html)
