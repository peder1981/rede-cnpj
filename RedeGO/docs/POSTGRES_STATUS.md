# 🐘 Status da Migração PostgreSQL

## ✅ CONCLUÍDO

### **1. PostgreSQL Configurado**
- ✅ Database `rede_cnpj` criado
- ✅ Usuário `rede_user` configurado
- ✅ 4 schemas criados (receita, rede, forensics, cache)
- ✅ 45 tabelas criadas (31 partições)
- ✅ 3 extensões instaladas (pg_trgm, btree_gin, unaccent)

**Conexão:**
```
Host: localhost
Porta: 5433
Database: rede_cnpj
Usuário: rede_user
Senha: rede_cnpj_2025
```

**String de Conexão:**
```
postgresql://rede_user:rede_cnpj_2025@localhost:5433/rede_cnpj?sslmode=disable
```

---

### **2. Código Atualizado**

#### **Arquivos Modificados:**

1. **`internal/config/config.go`**
   - ✅ Adicionado campo `PostgresURL`
   - ✅ Leitura do arquivo de configuração

2. **`internal/database/database.go`**
   - ✅ Suporte a PostgreSQL e SQLite
   - ✅ Funções helper: `GetDB()`, `IsPostgres()`
   - ✅ `AdaptQuery()` - Converte `?` para `$1, $2, ...`
   - ✅ `TablePrefix()` - Adiciona schema (receita., rede.)
   - ✅ Queries dos dicionários atualizadas

3. **`go.mod`**
   - ✅ Adicionado driver `github.com/lib/pq v1.10.9`

#### **Arquivos Criados:**

1. **`cmd/migrate/main.go`** - Migrador SQLite → PostgreSQL
2. **`rede.postgres.ini`** - Configuração exemplo
3. **`test_postgres.go`** - Script de teste
4. **`scripts/postgres/`** - Scripts SQL
   - `00_init_database.sql`
   - `01_schemas.sql`
   - `02_tables.sql`
   - `setup_postgres.sh`
   - `README.md`

5. **`Makefile`** - Novos comandos:
   - `make build-migrate` - Compila migrador
   - `make migrate` - Executa migração

---

### **3. Migração de Dados**

#### **Status Atual:**
```
🔄 EM ANDAMENTO

Empresas:         104.112 / ~50M  (0,2%)
Estabelecimentos: 0       / ~50M  (0%)
Sócios:           0       / ~26M  (0%)
```

**Comando em execução:**
```bash
./bin/rede-cnpj-migrate
```

**Progresso:** ~500 registros/segundo
**Tempo estimado:** 36-48 horas para migração completa

---

## 🔄 EM ANDAMENTO

### **Migração de Dados**
- ⏳ Migrando estabelecimentos (64M registros)
- ⏳ Aguardando sócios
- ⏳ Aguardando tabelas lookup
- ⏳ Aguardando rede/ligação

---

## 📋 PRÓXIMOS PASSOS

### **1. Aguardar Migração Completa**
Monitorar progresso:
```bash
# Ver logs da migração
tail -f /tmp/migrate.log  # (se configurado)

# Ou verificar contagens
psql -h localhost -p 5433 -U rede_user -d rede_cnpj -c "
SELECT 
  'empresas' as tabela, COUNT(*) FROM receita.empresas
UNION ALL
SELECT 'estabelecimento', COUNT(*) FROM receita.estabelecimento
UNION ALL
SELECT 'socios', COUNT(*) FROM receita.socios;
"
```

### **2. Atualizar Código das APIs**
Arquivos que precisam ser atualizados:
- [ ] `internal/api/handlers.go`
- [ ] `internal/api/cnpj.go`
- [ ] `internal/api/crossdata.go`
- [ ] `internal/api/forensics.go`
- [ ] `internal/graph/graph.go`

**Padrão de atualização:**
```go
// ANTES (SQLite)
db := database.GetDBReceita()
rows, err := db.Query("SELECT * FROM empresas WHERE cnpj_basico = ?", cnpj)

// DEPOIS (PostgreSQL compatível)
db := database.GetDB()
query := database.AdaptQuery("SELECT * FROM empresas WHERE cnpj_basico = ?")
query = fmt.Sprintf("SELECT * FROM %s WHERE cnpj_basico = $1", 
    database.TablePrefix("empresas"))
rows, err := db.Query(query, cnpj)
```

### **3. Criar Índices**
Após migração completa, criar índices para performance:
```bash
psql -h localhost -p 5433 -U rede_user -d rede_cnpj -f scripts/postgres/03_indexes.sql
```

**Índices prioritários:**
- B-Tree: `cnpj`, `cnpj_basico`, `cnpj_cpf_socio`
- GIN: Full-text search em `razao_social`, `nome_socio`
- GiST: Similaridade de texto
- BRIN: Datas (temporal)

### **4. Criar Views Materializadas**
Para queries frequentes:
```sql
-- Top sócios
CREATE MATERIALIZED VIEW cache.top_socios AS ...

-- Empresas ativas por UF
CREATE MATERIALIZED VIEW cache.empresas_por_uf AS ...

-- Refresh automático
REFRESH MATERIALIZED VIEW CONCURRENTLY cache.top_socios;
```

### **5. Testar Performance**
Comparar SQLite vs PostgreSQL:
```bash
# Benchmark de queries
go test -bench=. ./internal/api/...
```

### **6. Atualizar Configuração**
Editar `rede.ini`:
```ini
[BASE]
postgres_url = postgresql://rede_user:rede_cnpj_2025@localhost:5433/rede_cnpj?sslmode=disable

# Comentar SQLite (legacy)
# base_receita = bases/cnpj.db
# base_rede = bases/rede.db
# base_rede_search = bases/rede_search.db
```

### **7. Deploy**
```bash
# Compilar com PostgreSQL
make build

# Executar
./bin/rede-cnpj -conf_file=rede.postgres.ini
```

---

## 🧪 TESTES

### **Teste de Conexão**
```bash
go run test_postgres.go
```

**Resultado esperado:**
```
✅ Conectado ao PostgreSQL
✅ Empresas: 104112
✅ Estabelecimentos: 0 (aguardando migração)
✅ Sócios: 0 (aguardando migração)
```

### **Teste de Queries**
```bash
psql -h localhost -p 5433 -U rede_user -d rede_cnpj
```

```sql
-- Buscar empresa
SELECT * FROM receita.empresas WHERE cnpj_basico = '00000000';

-- Full-text search
SELECT razao_social 
FROM receita.empresas 
WHERE to_tsvector('portuguese', razao_social) @@ to_tsquery('portuguese', 'tecnologia');

-- Busca por similaridade
SELECT razao_social, similarity(razao_social, 'PETROBRAS') as sim
FROM receita.empresas
WHERE razao_social % 'PETROBRAS'
ORDER BY sim DESC
LIMIT 10;
```

---

## 📊 BENEFÍCIOS ESPERADOS

### **Performance**
- **33x mais rápido** em buscas por CPF/CNPJ
- **25x mais rápido** em full-text search
- **50x mais rápido** em agregações complexas
- **50x mais rápido** em queries de grafos

### **Recursos Avançados**
- ✅ Particionamento (28 UFs + 3 tipos de sócio)
- ✅ Full-text search nativo (português)
- ✅ Busca por similaridade (fuzzy)
- ✅ Índices avançados (GIN, GiST, BRIN)
- ✅ Views materializadas
- ✅ Paralelização de queries
- ✅ Suporte a JSON/JSONB

### **Escalabilidade**
- ✅ Suporta terabytes de dados
- ✅ Conexões concorrentes ilimitadas
- ✅ Replicação master-slave
- ✅ Backup incremental

---

## 🔧 TROUBLESHOOTING

### **Migração Lenta**
```bash
# Verificar configurações PostgreSQL
psql -h localhost -p 5433 -U rede_user -d rede_cnpj -c "
SHOW shared_buffers;
SHOW work_mem;
SHOW maintenance_work_mem;
"

# Ajustar se necessário
ALTER DATABASE rede_cnpj SET maintenance_work_mem = '2GB';
```

### **Erro de Conexão**
```bash
# Verificar se PostgreSQL está rodando
systemctl status postgresql

# Verificar porta
sudo netstat -tlnp | grep 5433

# Testar conexão
psql -h localhost -p 5433 -U rede_user -d rede_cnpj -c "SELECT 1;"
```

### **Rollback para SQLite**
Se necessário, basta comentar `postgres_url` no `rede.ini`:
```ini
[BASE]
# postgres_url = postgresql://...
base_receita = bases/cnpj.db
base_rede = bases/rede.db
```

---

## 📚 DOCUMENTAÇÃO

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [POSTGRES_MIGRATION.md](POSTGRES_MIGRATION.md) - Plano completo
- [scripts/postgres/README.md](scripts/postgres/README.md) - Scripts SQL

---

**Última atualização:** 2025-10-23 21:15
**Status:** 🟡 Migração em andamento
