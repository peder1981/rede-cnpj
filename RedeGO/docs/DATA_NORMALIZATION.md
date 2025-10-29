# Normalização e Validação de Dados - RedeCNPJ

## Visão Geral

Este documento descreve a estratégia de normalização e validação de dados implementada no processo de migração do SQLite para PostgreSQL. O objetivo é garantir a integridade dos dados e evitar erros de inserção causados por valores inválidos nos arquivos fonte da Receita Federal.

## Problema

Os dados fornecidos pela Receita Federal contêm diversos valores inválidos ou inconsistentes:

- **Datas inválidas**: `"0"`, `"00000000"`, strings vazias
- **Campos numéricos vazios**: `""`, `"0"`, `"00"`
- **Problemas de encoding**: caracteres ISO-8859-1 em vez de UTF-8
- **Espaços em branco**: valores com espaços desnecessários

Esses problemas causam erros durante a inserção no PostgreSQL:
```
pq: date/time field value out of range: "0"
pq: invalid input syntax for type numeric: ""
```

## Estratégia de Normalização

### 1. Sanitização de Strings (UTF-8)

**Função**: `sanitizeString(s string) string`

**Objetivo**: Corrigir problemas de encoding

**Processo**:
1. Verifica se a string já é UTF-8 válida
2. Tenta converter de ISO-8859-1 para UTF-8
3. Remove caracteres inválidos como último recurso

**Exemplo**:
```go
input:  "São Paulo" (com encoding incorreto)
output: "São Paulo" (UTF-8 válido)
```

### 2. Sanitização de Campos de Data

**Função**: `sanitizeDateString(ns sql.NullString) sql.NullString`

**Objetivo**: Normalizar e validar datas, convertendo valores inválidos em NULL

**Valores Inválidos Tratados**:
- Strings vazias: `""`
- Zeros: `"0"`, `"00"`, `"000"`, `"0000"`, `"00000000"`
- Anos inválidos: < 1900 ou > 2100
- Meses inválidos: < 01 ou > 12
- Dias inválidos: < 01 ou > 31

**Formatos Aceitos**:
- `YYYYMMDD` → convertido para `YYYY-MM-DD`
- `YYYY-MM-DD` → mantido como está

**Exemplos**:
```go
// Valores inválidos → NULL
"0"          → NULL
"00000000"   → NULL
""           → NULL
"18991231"   → NULL (ano < 1900)
"20231315"   → NULL (mês inválido)

// Valores válidos → formatados
"20231015"   → "2023-10-15"
"2023-10-15" → "2023-10-15"
```

### 3. Sanitização de Campos Numéricos

**Função**: `sanitizeNumericString(ns sql.NullString) sql.NullString`

**Objetivo**: Normalizar campos numéricos opcionais

**Valores Inválidos Tratados**:
- Strings vazias: `""`
- Zeros: `"0"`, `"00"`

**Exemplo**:
```go
"0"    → NULL
""     → NULL
"123"  → "123"
```

## Aplicação nas Tabelas

### Tabela: estabelecimento

**Campos de Data**:
- `data_situacao_cadastral` → `sanitizeDateString()`
- `data_inicio_atividades` → `sanitizeDateString()`
- `data_situacao_especial` → `sanitizeDateString()`

**Campos de String**:
- Todos os campos TEXT/VARCHAR → `sanitizeString()`

### Tabela: socios

**Campos de Data**:
- `data_entrada_sociedade` → `sanitizeDateString()`

**Campos de String**:
- Todos os campos TEXT/VARCHAR → `sanitizeString()`

### Tabela: simples

**Campos de Data**:
- `data_opcao_simples` → `sanitizeDateString()`
- `data_exclusao_simples` → `sanitizeDateString()`
- `data_opcao_mei` → `sanitizeDateString()`
- `data_exclusao_mei` → `sanitizeDateString()`

## Logs de Erro Melhorados

Quando ocorre um erro de inserção, o sistema agora registra:

```
⚠️  Erro ao inserir CNPJ 12345678000190: pq: date/time field value out of range: "0"
    Datas: situacao={0 true}, inicio={2023-10-15 true}, especial={0 false}
```

Isso facilita a identificação de:
- Qual CNPJ está causando o problema
- Quais valores de data estão sendo enviados
- Se o valor é válido ou NULL

## Testes Unitários

Os testes estão em `cmd/migrate/sanitize_test.go` e cobrem:

- ✅ Valores NULL
- ✅ Strings vazias
- ✅ Valores zero em diferentes formatos
- ✅ Datas válidas em diferentes formatos
- ✅ Datas com espaços em branco
- ✅ Anos, meses e dias inválidos
- ✅ Formatos inválidos

**Executar testes**:
```bash
cd cmd/migrate
go test -v
```

## Benefícios

1. **Integridade de Dados**: Garante que apenas valores válidos sejam inseridos
2. **Robustez**: Evita falhas na migração por dados inválidos
3. **Rastreabilidade**: Logs detalhados facilitam debug
4. **Manutenibilidade**: Código testado e documentado
5. **Performance**: Validação acontece em memória antes da inserção

## Próximos Passos

1. ✅ Implementar normalização de datas
2. ✅ Adicionar testes unitários
3. ✅ Melhorar logs de erro
4. 🔄 Monitorar migração e ajustar conforme necessário
5. 📊 Gerar relatório de dados normalizados vs. inválidos

## Estatísticas de Normalização

Durante a migração, o sistema registra:
- Total de registros processados
- Registros com datas normalizadas
- Registros com valores convertidos para NULL
- Erros de inserção (se houver)

Isso permite avaliar a qualidade dos dados fonte e identificar padrões de problemas.

## Referências

- [Documentação PostgreSQL - Date/Time Types](https://www.postgresql.org/docs/current/datatype-datetime.html)
- [Layout dos Dados Abertos CNPJ - Receita Federal](http://200.152.38.155/CNPJ/LAYOUT_DADOS_ABERTOS_CNPJ.pdf)
- [Go sql.NullString](https://pkg.go.dev/database/sql#NullString)
