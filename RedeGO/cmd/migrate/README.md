# Sistema de Migração RedeCNPJ

## Visão Geral

Sistema de migração de dados do SQLite para PostgreSQL com **normalização baseada em metadados de campos**.

## Arquitetura

```
cmd/migrate/
├── main.go              # Lógica principal de migração
├── normalizer.go        # Motor de normalização por tipo
├── schemas.go           # Metadados das tabelas
├── sanitize_test.go     # Testes de sanitização
├── normalizer_test.go   # Testes do normalizador
└── README.md           # Este arquivo
```

## Como Funciona

### 1. Definir Metadados da Tabela

```go
// schemas.go
func GetEstabelecimentoNormalizer() *Normalizer {
    n := NewNormalizer()
    
    // Registrar cada campo com seus metadados
    n.RegisterField(FieldMetadata{
        Name:      "cnpj",
        Type:      FieldTypeCNPJ,
        MaxLength: 14,
        Required:  true,
    })
    
    n.RegisterField(FieldMetadata{
        Name:     "data_situacao_cadastral",
        Type:     FieldTypeDate,
        Required: false,
    })
    
    return n
}
```

### 2. Usar na Migração

```go
// main.go
func migrateEstabelecimentos(src, dst *sql.DB) MigrationStats {
    // Criar normalizador específico
    normalizer := GetEstabelecimentoNormalizer()
    
    // Ler dados do SQLite
    rows, _ := src.Query("SELECT cnpj, data_situacao_cadastral FROM estabelecimento")
    
    for rows.Next() {
        var cnpj, data string
        rows.Scan(&cnpj, &data)
        
        // Normalizar campos
        cnpjNorm := normalizer.NormalizeString("cnpj", cnpj)
        dataNorm := normalizer.NormalizeNullString("data_situacao_cadastral", 
            sql.NullString{String: data, Valid: true})
        
        // Inserir no PostgreSQL
        stmt.Exec(cnpjNorm, dataNorm)
    }
}
```

## Tipos de Campo Suportados

### DATE - Datas

```go
n.RegisterField(FieldMetadata{
    Name:     "data_situacao_cadastral",
    Type:     FieldTypeDate,
    Required: false,
})

// Uso
dataNorm := normalizer.NormalizeNullString("data_situacao_cadastral", data)

// Exemplos
"20231015"   → "2023-10-15"
"0"          → NULL
"00000000"   → NULL
""           → NULL
```

### CNPJ - Cadastro Nacional de Pessoa Jurídica

```go
n.RegisterField(FieldMetadata{
    Name:      "cnpj",
    Type:      FieldTypeCNPJ,
    MaxLength: 14,
    Required:  true,
})

// Uso
cnpjNorm := normalizer.NormalizeString("cnpj", cnpj)

// Exemplos
"12.345.678/0001-90" → "12345678000190"
"12345678000190"     → "12345678000190"
"00000000000000"     → NULL
```

### CPF - Cadastro de Pessoa Física

```go
n.RegisterField(FieldMetadata{
    Name:      "representante_legal",
    Type:      FieldTypeCPF,
    MaxLength: 11,
    Required:  false,
})

// Uso
cpfNorm := normalizer.NormalizeString("representante_legal", cpf)

// Exemplos
"123.456.789-01" → "12345678901"
"12345678901"    → "12345678901"
"00000000000"    → NULL
```

### CEP - Código de Endereçamento Postal

```go
n.RegisterField(FieldMetadata{
    Name:      "cep",
    Type:      FieldTypeCEP,
    MaxLength: 8,
    Required:  false,
})

// Uso
cepNorm := normalizer.NormalizeString("cep", cep)

// Exemplos
"01310-100" → "01310100"
"01310100"  → "01310100"
"00000000"  → NULL
```

### EMAIL - Correio Eletrônico

```go
n.RegisterField(FieldMetadata{
    Name:     "correio_eletronico",
    Type:     FieldTypeEmail,
    Required: false,
})

// Uso
emailNorm := normalizer.NormalizeString("correio_eletronico", email)

// Exemplos
"TESTE@EXEMPLO.COM"   → "teste@exemplo.com"
"teste@exemplo.com"   → "teste@exemplo.com"
"teste.exemplo"       → NULL
```

### UF - Unidade Federativa

```go
n.RegisterField(FieldMetadata{
    Name:      "uf",
    Type:      FieldTypeUF,
    MaxLength: 2,
    Required:  true,
})

// Uso
ufNorm := normalizer.NormalizeString("uf", uf)

// Exemplos
"sp" → "SP"
"SP" → "SP"
"XX" → NULL
```

### PHONE - Telefone

```go
n.RegisterField(FieldMetadata{
    Name:      "telefone1",
    Type:      FieldTypePhone,
    MaxLength: 8,
    Required:  false,
})

// Uso
phoneNorm := normalizer.NormalizeString("telefone1", phone)

// Exemplos
"(11) 98765-4321" → "11987654321"
"11987654321"     → "11987654321"
"0"               → NULL
```

### CODE - Códigos Genéricos

```go
n.RegisterField(FieldMetadata{
    Name:      "cnae_fiscal",
    Type:      FieldTypeCode,
    MaxLength: 7,
    Required:  false,
})

// Uso
codeNorm := normalizer.NormalizeString("cnae_fiscal", code)

// Exemplos
"  1234567  " → "1234567"
"0"           → NULL
""            → NULL
```

## Executar Migração

### Pré-requisitos

```bash
# PostgreSQL rodando
docker-compose up -d postgres

# SQLite com dados
ls bases/cnpj.db
```

### Executar

```bash
cd cmd/migrate
go run main.go normalizer.go schemas.go
```

### Saída Esperada

```
╔════════════════════════════════════════════════════════════════╗
║                                                                ║
║     🔄 RedeCNPJ - Migração SQLite → PostgreSQL                ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝

📂 Conectando ao SQLite...
🐘 Conectando ao PostgreSQL...
✅ Conexões estabelecidas!

📊 Migrando tabela: empresas
   Total de registros: 64888615
   Progresso: 10000/64888615 (0.0%)
   Progresso: 20000/64888615 (0.0%)
   ...
✅ Empresas migradas: 64888615 em 3h12m11s

📊 Migrando tabela: estabelecimento
   Total de registros: 68048886
   Progresso: 10000/68048886 (0.0%)
   Progresso: 20000/68048886 (0.0%)
   ...
✅ Estabelecimentos migrados: 68048886 em 4h30m22s

╔════════════════════════════════════════════════════════════════╗
║                                                                ║
║     ✅ MIGRAÇÃO CONCLUÍDA COM SUCESSO!                        ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝
```

## Testes

### Executar Todos os Testes

```bash
cd cmd/migrate
go test -v
```

### Executar Teste Específico

```bash
go test -v -run TestNormalizerCNPJ
```

### Cobertura de Testes

```bash
go test -cover
```

## Adicionar Nova Tabela

### 1. Criar Normalizador

```go
// schemas.go
func GetMinhaNovaTabela() *Normalizer {
    n := NewNormalizer()
    
    n.RegisterField(FieldMetadata{
        Name:      "campo1",
        Type:      FieldTypeText,
        Required:  true,
    })
    
    n.RegisterField(FieldMetadata{
        Name:     "campo2",
        Type:     FieldTypeDate,
        Required: false,
    })
    
    return n
}
```

### 2. Criar Função de Migração

```go
// main.go
func migrateMinhaNovaTabela(src, dst *sql.DB) MigrationStats {
    normalizer := GetMinhaNovaTabela()
    
    // ... lógica de migração
    
    campo1Norm := normalizer.NormalizeString("campo1", campo1)
    campo2Norm := normalizer.NormalizeNullString("campo2", campo2)
    
    stmt.Exec(campo1Norm, campo2Norm)
}
```

### 3. Adicionar no main()

```go
func main() {
    // ...
    
    stat = migrateMinhaNovaTabela(srcDB, dstDB)
    stats = append(stats, stat)
    
    // ...
}
```

## Adicionar Novo Tipo de Campo

### 1. Definir Tipo

```go
// normalizer.go
const (
    // ... tipos existentes
    FieldTypeMeuNovoTipo FieldType = "MEU_NOVO_TIPO"
)
```

### 2. Implementar Normalização

```go
// normalizer.go
func (n *Normalizer) normalizeMeuNovoTipo(value string) sql.NullString {
    // Lógica de normalização
    if value == "" {
        return sql.NullString{Valid: false}
    }
    
    // Transformações
    normalized := strings.ToUpper(value)
    
    return sql.NullString{String: normalized, Valid: true}
}
```

### 3. Adicionar no Switch

```go
// normalizer.go
func (n *Normalizer) NormalizeString(fieldName string, value string) sql.NullString {
    switch meta.Type {
    // ... casos existentes
    case FieldTypeMeuNovoTipo:
        return n.normalizeMeuNovoTipo(trimmed)
    }
}
```

### 4. Criar Testes

```go
// normalizer_test.go
func TestNormalizerMeuNovoTipo(t *testing.T) {
    n := NewNormalizer()
    n.RegisterField(FieldMetadata{
        Name: "campo_teste",
        Type: FieldTypeMeuNovoTipo,
    })
    
    result := n.NormalizeString("campo_teste", "valor")
    
    if result.String != "VALOR" {
        t.Errorf("Esperado 'VALOR', obtido %q", result.String)
    }
}
```

## Troubleshooting

### Erro: "date/time field value out of range"

**Causa**: Data inválida não foi normalizada  
**Solução**: Verificar se campo está registrado com `FieldTypeDate`

```go
n.RegisterField(FieldMetadata{
    Name:     "data_campo",
    Type:     FieldTypeDate,  // ← Importante!
    Required: false,
})
```

### Erro: "invalid input syntax for type numeric"

**Causa**: Campo numérico com valor inválido  
**Solução**: Usar `NormalizeFloat64` para campos numéricos

```go
capitalNorm := normalizer.NormalizeFloat64("capital_social", capital)
```

### Erro: "declared and not used"

**Causa**: Variável normalizada não está sendo usada  
**Solução**: Usar variável no `stmt.Exec()`

```go
cnpjNorm := normalizer.NormalizeString("cnpj", cnpj)
stmt.Exec(cnpjNorm, ...)  // ← Usar cnpjNorm, não cnpj
```

## Performance

### Otimizações Implementadas

- ✅ Batch processing (10.000 registros por transação)
- ✅ Prepared statements (reutilizados)
- ✅ Validação em memória (sem queries extras)
- ✅ Índices no PostgreSQL (criados antes da migração)

### Tempos Esperados

| Tabela | Registros | Tempo Estimado |
|--------|-----------|----------------|
| empresas | 64.8M | ~3h |
| estabelecimento | 68.0M | ~4.5h |
| socios | 50.0M | ~3h |
| simples | 20.0M | ~1h |
| **Total** | **~200M** | **~11.5h** |

## Documentação Adicional

- [DATA_NORMALIZATION.md](../../docs/DATA_NORMALIZATION.md) - Estratégia de normalização
- [FIELD_BASED_NORMALIZATION.md](../../docs/FIELD_BASED_NORMALIZATION.md) - Sistema baseado em metadados
- [NORMALIZATION_SUMMARY.md](../../docs/NORMALIZATION_SUMMARY.md) - Resumo executivo
- [MIGRATION_IMPROVEMENTS.md](../../docs/MIGRATION_IMPROVEMENTS.md) - Melhorias implementadas

## Suporte

Para dúvidas ou problemas:
1. Verificar logs de erro
2. Consultar documentação
3. Executar testes unitários
4. Abrir issue no repositório
