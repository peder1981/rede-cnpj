# 🚀 Importação Otimizada - Guia Rápido

## ⚡ Problema Resolvido

**ANTES**: Importação levava 4+ dias
- Passo 1: Importar para SQLite (~2 dias)
- Passo 2: Migrar para PostgreSQL (~2+ dias)
- **Total**: ~4+ dias ❌

**AGORA**: Importação leva ~2 dias
- Passo único: Importar diretamente para PostgreSQL
- **Total**: ~2 dias ✅

## 📋 Pré-requisitos

1. **PostgreSQL instalado e rodando**
```bash
# Verificar se PostgreSQL está rodando
sudo systemctl status postgresql

# Criar banco de dados
sudo -u postgres psql -c "CREATE DATABASE rede_cnpj;"
sudo -u postgres psql -c "CREATE USER rede_user WITH PASSWORD 'rede_cnpj_2025';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE rede_cnpj TO rede_user;"
```

2. **Arquivo de configuração**

Crie ou edite `rede.postgres.ini`:
```ini
[BASE]
postgres_url = postgresql://rede_user:rede_cnpj_2025@localhost:5432/rede_cnpj?sslmode=disable
base_receita = bases/cnpj.db
base_rede = bases/rede.db
base_rede_search = bases/rede_search.db
pasta_arquivos = arquivos

[ETC]
limiter_padrao = 20/minute
limite_registros_camada = 1000
tempo_maximo_consulta = 10.0
```

## 🎯 Como Executar

### Opção 1: Importação Completa (Recomendado)

```bash
cd RedeGO

# Executar importação completa diretamente para PostgreSQL
go run cmd/importer/main.go -config rede.postgres.ini -all
```

Isso irá:
1. ✅ Baixar arquivos ZIP da Receita Federal
2. ✅ Processar CSVs com normalização
3. ✅ Importar diretamente para PostgreSQL
4. ✅ Criar índices otimizados
5. ✅ Criar tabelas de ligação
6. ✅ Criar índices de busca

### Opção 2: Passo a Passo

```bash
# 1. Apenas baixar arquivos
go run cmd/importer/main.go -config rede.postgres.ini -download

# 2. Processar e importar
go run cmd/importer/main.go -config rede.postgres.ini -process

# 3. Criar tabelas de ligação
go run cmd/importer/main.go -config rede.postgres.ini -links

# 4. Criar índices de busca
go run cmd/importer/main.go -config rede.postgres.ini -search
```

## 📊 Monitoramento

Durante a importação, você verá:

```
╔════════════════════════════════════════════════════════════════╗
║                                                                ║
║     📥 RedeCNPJ - Importador de Dados da Receita Federal      ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝

ℹ️  Carregando configuração de rede.postgres.ini...
🐘 Banco de dados: PostgreSQL (importação direta)
✅ Usando PostgreSQL para importação direta
✅ Schemas PostgreSQL criados/verificados

📦 Iniciando processamento dos arquivos...
📋 Encontrados 45 arquivos ZIP para processar

[1/45] 📦 Processando Empresas0.zip...
    Importando EMPRE.D81234 para tabela empresas...
      100000 registros...
      200000 registros...
      ✅ 534123 registros importados

[2/45] 📦 Processando Estabelecimentos0.zip...
    Importando ESTABELE.D81234 para tabela estabelecimento...
      100000 registros...
      200000 registros...
      ...
```

## 🔍 Verificação

Após a importação, verifique os dados:

```bash
# Conectar ao PostgreSQL
psql -U rede_user -d rede_cnpj

# Verificar contagem de registros
SELECT 'empresas' as tabela, COUNT(*) FROM receita.empresas
UNION ALL
SELECT 'estabelecimento', COUNT(*) FROM receita.estabelecimento
UNION ALL
SELECT 'socios', COUNT(*) FROM receita.socios
UNION ALL
SELECT 'simples', COUNT(*) FROM receita.simples;

# Verificar alguns registros
SELECT * FROM receita.empresas LIMIT 5;
SELECT * FROM receita.estabelecimento LIMIT 5;
```

## ⚙️ Configurações Avançadas

### Ajustar Performance

No arquivo `processor.go`, você pode ajustar:

```go
const (
    batchSize = 10000  // Tamanho do lote (padrão: 10000)
)
```

### Configurações PostgreSQL

Para melhor performance durante importação, ajuste temporariamente:

```sql
-- Desabilitar autovacuum durante importação
ALTER TABLE receita.estabelecimento SET (autovacuum_enabled = false);

-- Aumentar work_mem
SET work_mem = '256MB';

-- Desabilitar fsync (CUIDADO: apenas durante importação)
ALTER SYSTEM SET fsync = off;
SELECT pg_reload_conf();

-- Após importação, reverter
ALTER SYSTEM SET fsync = on;
SELECT pg_reload_conf();
ALTER TABLE receita.estabelecimento SET (autovacuum_enabled = true);
```

## 🐛 Troubleshooting

### Erro: "connection refused"
```bash
# Verificar se PostgreSQL está rodando
sudo systemctl status postgresql
sudo systemctl start postgresql
```

### Erro: "permission denied"
```bash
# Verificar permissões do usuário
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE rede_cnpj TO rede_user;"
sudo -u postgres psql -d rede_cnpj -c "GRANT ALL ON SCHEMA receita TO rede_user;"
sudo -u postgres psql -d rede_cnpj -c "GRANT ALL ON SCHEMA rede TO rede_user;"
```

### Erro: "out of memory"
```bash
# Aumentar memória disponível no PostgreSQL
sudo vim /etc/postgresql/*/main/postgresql.conf

# Ajustar:
shared_buffers = 4GB
work_mem = 256MB
maintenance_work_mem = 1GB
effective_cache_size = 12GB

# Reiniciar PostgreSQL
sudo systemctl restart postgresql
```

### Importação muito lenta
1. Verificar se há outros processos consumindo recursos
2. Aumentar `batchSize` no código
3. Desabilitar autovacuum temporariamente
4. Usar SSD ao invés de HDD
5. Aumentar `shared_buffers` do PostgreSQL

## 📈 Métricas Esperadas

Com hardware moderno (SSD, 16GB RAM, 8 cores):

| Tabela           | Registros    | Tempo Estimado |
|------------------|--------------|----------------|
| Empresas         | ~50M         | ~2h            |
| Estabelecimento  | ~55M         | ~8h            |
| Sócios           | ~25M         | ~4h            |
| Simples          | ~20M         | ~2h            |
| Tabelas lookup   | ~10K         | ~1min          |
| **TOTAL**        | **~150M**    | **~16-20h**    |

## 🎓 Diferenças vs Método Antigo

| Aspecto              | Método Antigo        | Método Novo          |
|----------------------|----------------------|----------------------|
| Banco intermediário  | SQLite               | Nenhum               |
| Etapas               | 2 (import + migrate) | 1 (import direto)    |
| Normalização         | Após importação      | Durante importação   |
| Validação            | Após importação      | Durante importação   |
| Tempo total          | ~4+ dias             | ~2 dias              |
| Uso de disco         | 2x (SQLite + PG)     | 1x (apenas PG)       |
| Risco de erro        | Alto (2 etapas)      | Baixo (1 etapa)      |

## ✅ Checklist de Importação

- [ ] PostgreSQL instalado e configurado
- [ ] Banco de dados `rede_cnpj` criado
- [ ] Usuário `rede_user` com permissões
- [ ] Arquivo `rede.postgres.ini` configurado
- [ ] Espaço em disco suficiente (~200GB)
- [ ] Memória RAM suficiente (16GB+)
- [ ] Conexão de internet estável (para download)
- [ ] Executar importação: `go run cmd/importer/main.go -config rede.postgres.ini -all`
- [ ] Aguardar conclusão (~2 dias)
- [ ] Verificar dados no PostgreSQL
- [ ] Criar backup do banco

## 📚 Documentação Adicional

- [IMPORTER_DIRECT_POSTGRES.md](docs/IMPORTER_DIRECT_POSTGRES.md) - Detalhes técnicos
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) - Arquitetura do sistema
- [POSTGRES_MIGRATION.md](docs/POSTGRES_MIGRATION.md) - Migração PostgreSQL

## 🆘 Suporte

Em caso de problemas:
1. Verificar logs do PostgreSQL: `/var/log/postgresql/`
2. Verificar logs da aplicação
3. Abrir issue no GitHub com detalhes do erro
