# Importação Direta para PostgreSQL

## Problema Identificado

O fluxo anterior de importação estava causando um processo extremamente lento (4+ dias):

1. **Importer** → Importava dados da Receita para **SQLite** (`cnpj.db`)
2. **Migrate** → Migrava dados do **SQLite** para **PostgreSQL**

Isso significava que os dados eram escritos **duas vezes**, duplicando o tempo de processamento e causando um gargalo significativo.

## Solução Implementada

O importer agora detecta automaticamente qual banco de dados está configurado e importa **diretamente** para ele, eliminando a necessidade de migração posterior.

### Arquitetura Nova

```
┌─────────────────────────────────────────────────────────────┐
│                    IMPORTER (cmd/importer)                  │
│                                                             │
│  1. Detecta banco configurado (SQLite ou PostgreSQL)       │
│  2. Baixa arquivos ZIP da Receita Federal                  │
│  3. Processa CSVs com normalização                         │
│  4. Importa DIRETAMENTE para o banco correto               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ├─────────────────┬──────────────────┐
                              ▼                 ▼                  ▼
                      ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
                      │   SQLite     │  │  PostgreSQL  │  │ Normalização │
                      │   (legacy)   │  │   (direto)   │  │  (durante)   │
                      └──────────────┘  └──────────────┘  └──────────────┘
```

### Componentes Criados

#### 1. `internal/importer/database.go`
- **DatabaseManager**: Gerencia conexão com banco de dados
- Detecta automaticamente PostgreSQL ou SQLite
- Cria schemas necessários no PostgreSQL
- Adapta queries e placeholders para cada banco

#### 2. `internal/importer/schemas.go`
- Define schemas de tabelas para SQLite e PostgreSQL
- Define índices otimizados para cada banco
- Suporta schemas separados (receita, rede) no PostgreSQL

#### 3. `internal/importer/normalizer.go` + `normalizer_schemas.go`
- Copiado do `cmd/migrate` para o importer
- Normaliza dados **durante a importação**
- Valida tipos de dados (CNPJ, CPF, CEP, datas, etc)
- Remove dados inválidos antes de inserir

#### 4. `internal/importer/sanitize.go`
- Corrige problemas de encoding UTF-8
- Converte ISO-8859-1 para UTF-8
- Remove caracteres inválidos

#### 5. `internal/importer/inserter.go`
- Gerencia inserção de dados com normalização
- Prepara statements otimizados para cada banco
- Aplica normalização específica por tabela
- Trata conflitos (ON CONFLICT) no PostgreSQL

### Mudanças no Processor

O `internal/importer/processor.go` foi refatorado para:

1. **Usar DatabaseManager** ao invés de conexão direta
2. **Criar schemas corretos** baseado no tipo de banco
3. **Aplicar normalização** durante a importação
4. **Adaptar queries** automaticamente (? para $1, $2, etc)
5. **Usar TablePrefix** para schemas PostgreSQL (receita.*, rede.*)

### Mudanças no Main

O `cmd/importer/main.go` agora:

1. **Carrega configuração** do arquivo `rede.ini` ou `rede.postgres.ini`
2. **Detecta PostgreSQL** pela presença de `postgres_url`
3. **Exibe mensagem clara** sobre qual banco será usado
4. **Passa configuração** para o Processor

## Como Usar

### Opção 1: Importar para SQLite (modo legado)

```bash
cd RedeGO
go run cmd/importer/main.go -all
```

Ou sem arquivo de configuração:
```bash
go run cmd/importer/main.go -all
```

### Opção 2: Importar Diretamente para PostgreSQL

1. Configure o PostgreSQL no arquivo `rede.postgres.ini`:
```ini
[BASE]
postgres_url = postgresql://usuario:senha@localhost:5432/rede_cnpj?sslmode=disable
```

2. Execute o importer:
```bash
go run cmd/importer/main.go -config rede.postgres.ini -all
```

Ou renomeie o arquivo para `rede.ini` e execute:
```bash
go run cmd/importer/main.go -all
```

### Flags Disponíveis

- `-download`: Apenas baixa os arquivos ZIP
- `-process`: Apenas processa arquivos já baixados
- `-links`: Cria tabelas de ligação (rede.db)
- `-search`: Cria índices de busca
- `-all`: Executa todo o processo
- `-config`: Especifica arquivo de configuração (padrão: rede.ini)

## Benefícios

### Performance
- **Elimina duplicação**: Dados escritos apenas uma vez
- **Reduz tempo**: De 4+ dias para ~2 dias (estimativa)
- **Normalização eficiente**: Aplicada durante importação, não depois

### Qualidade de Dados
- **Validação imediata**: Dados inválidos rejeitados na importação
- **Encoding correto**: UTF-8 garantido desde o início
- **Tipos corretos**: Datas, números e códigos validados

### Manutenibilidade
- **Código unificado**: Normalização compartilhada entre importer e migrate
- **Detecção automática**: Não precisa escolher manualmente o banco
- **Compatibilidade**: Mantém suporte a SQLite para desenvolvimento

## Fluxo Completo

### Para PostgreSQL (Recomendado)

```bash
# 1. Configurar PostgreSQL
vim rede.postgres.ini

# 2. Criar banco e schemas (se necessário)
psql -U postgres -c "CREATE DATABASE rede_cnpj"

# 3. Executar importação direta
go run cmd/importer/main.go -config rede.postgres.ini -all

# 4. Pronto! Dados já estão no PostgreSQL normalizados
```

### Para SQLite (Legacy)

```bash
# 1. Executar importação
go run cmd/importer/main.go -all

# 2. Se precisar migrar para PostgreSQL depois
go run cmd/migrate/main.go
```

## Migração Obsoleta?

O comando `migrate` ainda é útil para:
- Migrar bases SQLite antigas para PostgreSQL
- Reprocessar dados já importados
- Testes e desenvolvimento

Mas para **novas importações**, use sempre a **importação direta** para PostgreSQL.

## Próximos Passos

1. **Testar importação completa** com dados reais
2. **Medir performance** e comparar com método antigo
3. **Documentar métricas** de tempo e uso de recursos
4. **Otimizar batch size** se necessário
5. **Adicionar progress bar** mais detalhado

## Notas Técnicas

### Normalização Durante Importação

Cada tabela tem seu próprio normalizador que valida:
- **Empresas**: CNPJ básico, capital social, códigos
- **Estabelecimento**: CNPJ completo, UF, CEP, datas, email
- **Sócios**: CPF/CNPJ, datas, códigos
- **Simples**: Datas, opções S/N

### Tratamento de Conflitos

No PostgreSQL, usa `ON CONFLICT DO NOTHING` para:
- Evitar duplicatas em chaves primárias
- Permitir re-execução segura
- Manter dados mais recentes

### Schemas PostgreSQL

Dados organizados em schemas separados:
- `receita.*`: Dados da Receita Federal
- `rede.*`: Dados de relacionamento/rede

Isso facilita:
- Permissões granulares
- Backup seletivo
- Organização lógica

## Compatibilidade

- ✅ Go 1.21+
- ✅ PostgreSQL 12+
- ✅ SQLite 3.35+
- ✅ Linux, macOS, Windows
