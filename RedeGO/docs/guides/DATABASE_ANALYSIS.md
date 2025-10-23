# 🗄️ Análise Completa das Estruturas de Banco de Dados

## Bancos de Dados Disponíveis

### 1. **cnpj.db** - Dados da Receita Federal

#### Tabelas de Códigos (Lookup)
- **cnae** - Classificação Nacional de Atividades Econômicas
  - `codigo TEXT` - Código CNAE
  - `descricao TEXT` - Descrição da atividade

- **motivo** - Motivos de situação cadastral
  - `codigo TEXT` - Código do motivo
  - `descricao TEXT` - Descrição (baixa, suspensão, etc)

- **municipio** - Municípios brasileiros
  - `codigo TEXT` - Código IBGE
  - `descricao TEXT` - Nome do município

- **natureza_juridica** - Natureza jurídica das empresas
  - `codigo TEXT` - Código da natureza
  - `descricao TEXT` - Descrição (SA, LTDA, MEI, etc)

- **pais** - Países
  - `codigo TEXT` - Código do país
  - `descricao TEXT` - Nome do país

- **qualificacao_socio** - Qualificação de sócios
  - `codigo TEXT` - Código da qualificação
  - `descricao TEXT` - Descrição (Administrador, Sócio, etc)

#### Tabelas Principais

- **empresas** - Dados cadastrais básicos
  - `cnpj_basico TEXT` - 8 primeiros dígitos do CNPJ
  - `razao_social TEXT` - **NOME COMPLETO DA EMPRESA**
  - `natureza_juridica TEXT` - Tipo jurídico
  - `qualificacao_responsavel TEXT` - Qualificação do responsável
  - `capital_social REAL` - Capital social em R$
  - `porte_empresa TEXT` - Porte (MEI, ME, EPP, etc)
  - `ente_federativo_responsavel TEXT` - Para empresas públicas

- **estabelecimento** - Dados de cada estabelecimento (matriz/filial)
  - `cnpj_basico TEXT` - 8 primeiros dígitos
  - `cnpj_ordem TEXT` - 4 dígitos da ordem
  - `cnpj_dv TEXT` - 2 dígitos verificadores
  - `cnpj TEXT` - **CNPJ COMPLETO (14 dígitos)**
  - `matriz_filial TEXT` - '1'=Matriz, '2'=Filial
  - `nome_fantasia TEXT` - **NOME FANTASIA**
  - `situacao_cadastral TEXT` - Ativa, Baixada, Suspensa, etc
  - `data_situacao_cadastral TEXT` - Data da situação
  - `motivo_situacao_cadastral TEXT` - Motivo da baixa/suspensão
  - `nome_cidade_exterior TEXT` - Para empresas no exterior
  - `pais TEXT` - Código do país
  - `data_inicio_atividades TEXT` - **DATA DE ABERTURA**
  - `cnae_fiscal TEXT` - **ATIVIDADE PRINCIPAL**
  - `cnae_fiscal_secundaria TEXT` - **ATIVIDADES SECUNDÁRIAS**
  - `tipo_logradouro TEXT` - Rua, Av, etc
  - `logradouro TEXT` - **ENDEREÇO**
  - `numero TEXT` - Número
  - `complemento TEXT` - Complemento
  - `bairro TEXT` - Bairro
  - `cep TEXT` - **CEP**
  - `uf TEXT` - **ESTADO**
  - `municipio TEXT` - Código do município
  - `ddd1 TEXT` - **DDD TELEFONE 1**
  - `telefone1 TEXT` - **TELEFONE 1**
  - `ddd2 TEXT` - DDD telefone 2
  - `telefone2 TEXT` - Telefone 2
  - `ddd_fax TEXT` - DDD fax
  - `fax TEXT` - Fax
  - `correio_eletronico TEXT` - **EMAIL**
  - `situacao_especial TEXT` - Situação especial
  - `data_situacao_especial TEXT` - Data situação especial

- **socios** - **SÓCIOS E ADMINISTRADORES**
  - `cnpj TEXT` - **CNPJ DA EMPRESA (14 dígitos)**
  - `cnpj_basico TEXT` - 8 primeiros dígitos
  - `identificador_de_socio TEXT` - 1=PF, 2=PJ, 3=Estrangeiro
  - `nome_socio TEXT` - **NOME COMPLETO DO SÓCIO**
  - `cnpj_cpf_socio TEXT` - **CPF OU CNPJ DO SÓCIO (SEM MÁSCARA)**
  - `qualificacao_socio TEXT` - **CARGO/FUNÇÃO**
  - `data_entrada_sociedade TEXT` - **DATA DE ENTRADA**
  - `pais TEXT` - País (para estrangeiros)
  - `representante_legal TEXT` - **CPF DO REPRESENTANTE LEGAL**
  - `nome_representante TEXT` - **NOME DO REPRESENTANTE**
  - `qualificacao_representante_legal TEXT` - Qualificação do representante
  - `faixa_etaria TEXT` - **FAIXA ETÁRIA DO SÓCIO**

- **simples** - Opção pelo Simples Nacional e MEI
  - `cnpj_basico TEXT` - 8 primeiros dígitos
  - `opcao_simples TEXT` - S=Sim, N=Não
  - `data_opcao_simples TEXT` - Data de opção
  - `data_exclusao_simples TEXT` - Data de exclusão
  - `opcao_mei TEXT` - **S=MEI, N=Não**
  - `data_opcao_mei TEXT` - Data de opção MEI
  - `data_exclusao_mei TEXT` - Data de exclusão MEI

### 2. **rede.db** - Rede de Relacionamentos

- **ligacao** - Relacionamentos entre entidades
  - `id1 TEXT` - **ID ORIGEM** (PJ_cnpj, PF_cpf-nome, PE_nome)
  - `id2 TEXT` - **ID DESTINO**
  - `descricao TEXT` - **TIPO DE RELACIONAMENTO**
  - `cnpj TEXT` - CNPJ relacionado
  - `peso INTEGER` - Peso da ligação

### 3. **rede_search.db** - Índice de Busca FTS5

- **id_search** - Índice full-text search
  - `id_descricao TEXT` - **ID + NOME** para busca rápida

## 🎯 Campos Censurados no Código Atual

### Localizado em: `rede_sqlite_cnpj.py` linhas 1431-1437

```python
if d['natureza_juridica'] in ('2135', '4120'): #remove empresario individual, produtor rural
    ts = '#INFORMAÇÃO EDITADA#'
    d['endereco'] = ts
    d['telefone1'] = ts
    d['telefone2'] = ''
    d['fax'] = ''
    d['correio_eletronico'] = ts
    d['cep'] = ts
```

**Naturezas Jurídicas Censuradas:**
- **2135** - Empresário Individual
- **4120** - Produtor Rural (Pessoa Física)

**Campos Censurados:**
- Endereço completo
- Telefones (1 e 2)
- Fax
- Email
- CEP

## 🔓 Dados Disponíveis SEM Censura

### Dados Pessoais Completos
1. **CPF completo** (sem máscara) - `cnpj_cpf_socio`
2. **Nome completo** - `nome_socio`
3. **Faixa etária** - `faixa_etaria`
4. **Data de entrada na sociedade** - `data_entrada_sociedade`
5. **Cargo/Qualificação** - `qualificacao_socio`
6. **Representante legal** (CPF + Nome)

### Dados Empresariais Completos
1. **CNPJ completo** (14 dígitos)
2. **Razão social completa**
3. **Nome fantasia**
4. **Endereço completo** (rua, número, complemento, bairro, CEP, cidade, UF)
5. **Telefones** (até 2 + fax)
6. **Email**
7. **Capital social**
8. **Data de abertura**
9. **Situação cadastral**
10. **CNAEs** (principal + secundárias)

## 🔗 Possibilidades de Cruzamento de Dados

### 1. **Cruzamento CPF → Empresas**
```sql
SELECT DISTINCT 
    s.cnpj_cpf_socio as cpf,
    s.nome_socio,
    e.cnpj,
    e.razao_social,
    est.nome_fantasia,
    s.qualificacao_socio,
    s.data_entrada_sociedade
FROM socios s
JOIN estabelecimento est ON s.cnpj = est.cnpj
JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
WHERE s.cnpj_cpf_socio = ?
ORDER BY s.data_entrada_sociedade
```

**Resultado:** Todas as empresas onde uma pessoa é sócia/administradora

### 2. **Cruzamento CNPJ → Sócios**
```sql
SELECT 
    s.nome_socio,
    s.cnpj_cpf_socio,
    s.qualificacao_socio,
    s.faixa_etaria,
    s.data_entrada_sociedade,
    s.representante_legal,
    s.nome_representante
FROM socios s
WHERE s.cnpj = ?
ORDER BY s.qualificacao_socio, s.nome_socio
```

**Resultado:** Todos os sócios de uma empresa

### 3. **Cruzamento Sócios em Comum**
```sql
SELECT DISTINCT
    s1.cnpj as cnpj1,
    s2.cnpj as cnpj2,
    s1.nome_socio,
    s1.cnpj_cpf_socio,
    s1.qualificacao_socio as qualif_empresa1,
    s2.qualificacao_socio as qualif_empresa2
FROM socios s1
JOIN socios s2 ON s1.cnpj_cpf_socio = s2.cnpj_cpf_socio
WHERE s1.cnpj = ? AND s2.cnpj = ? AND s1.cnpj != s2.cnpj
```

**Resultado:** Pessoas que são sócias de ambas as empresas

### 4. **Rede de Empresas de uma Pessoa**
```sql
WITH empresas_pessoa AS (
    SELECT DISTINCT cnpj, qualificacao_socio
    FROM socios
    WHERE cnpj_cpf_socio = ?
)
SELECT 
    ep.cnpj,
    e.razao_social,
    est.nome_fantasia,
    ep.qualificacao_socio,
    s2.nome_socio as outros_socios,
    s2.cnpj_cpf_socio as cpf_outros_socios
FROM empresas_pessoa ep
JOIN estabelecimento est ON ep.cnpj = est.cnpj
JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
LEFT JOIN socios s2 ON ep.cnpj = s2.cnpj AND s2.cnpj_cpf_socio != ?
ORDER BY ep.cnpj, s2.nome_socio
```

**Resultado:** Todas as empresas de uma pessoa + outros sócios dessas empresas

### 5. **Empresas no Mesmo Endereço**
```sql
SELECT 
    est1.cnpj,
    e1.razao_social,
    est1.nome_fantasia,
    est1.situacao_cadastral
FROM estabelecimento est1
JOIN empresas e1 ON est1.cnpj_basico = e1.cnpj_basico
WHERE est1.cep = ? 
  AND est1.logradouro = ?
  AND est1.numero = ?
ORDER BY est1.razao_social
```

**Resultado:** Empresas que compartilham o mesmo endereço

### 6. **Empresas com Mesmo Email/Telefone**
```sql
SELECT 
    est.cnpj,
    e.razao_social,
    est.nome_fantasia,
    est.correio_eletronico,
    est.telefone1
FROM estabelecimento est
JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
WHERE est.correio_eletronico = ? OR est.telefone1 = ?
```

**Resultado:** Empresas que compartilham email ou telefone

### 7. **Representantes Legais**
```sql
SELECT 
    s.cnpj,
    e.razao_social,
    s.nome_socio as socio_menor,
    s.cnpj_cpf_socio as cpf_menor,
    s.representante_legal as cpf_representante,
    s.nome_representante
FROM socios s
JOIN estabelecimento est ON s.cnpj = est.cnpj
JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
WHERE s.representante_legal IS NOT NULL AND s.representante_legal != ''
```

**Resultado:** Menores de idade com representantes legais

### 8. **Empresas Estrangeiras**
```sql
SELECT 
    est.cnpj,
    e.razao_social,
    est.nome_fantasia,
    est.nome_cidade_exterior,
    p.descricao as pais,
    est.correio_eletronico
FROM estabelecimento est
JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
LEFT JOIN pais p ON est.pais = p.codigo
WHERE est.uf = 'EX'
```

**Resultado:** Empresas com sede no exterior

### 9. **Sócios Estrangeiros**
```sql
SELECT 
    s.cnpj,
    e.razao_social,
    s.nome_socio,
    p.descricao as pais,
    s.qualificacao_socio
FROM socios s
JOIN estabelecimento est ON s.cnpj = est.cnpj
JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
LEFT JOIN pais p ON s.pais = p.codigo
WHERE s.identificador_de_socio = '3'
```

**Resultado:** Sócios estrangeiros

### 10. **Empresas por Faixa Etária dos Sócios**
```sql
SELECT 
    s.faixa_etaria,
    COUNT(DISTINCT s.cnpj) as total_empresas,
    COUNT(DISTINCT s.cnpj_cpf_socio) as total_socios
FROM socios s
WHERE s.faixa_etaria IS NOT NULL
GROUP BY s.faixa_etaria
ORDER BY s.faixa_etaria
```

**Resultado:** Distribuição de empresas por faixa etária dos sócios

### 11. **Empresas Baixadas com Sócios Ativos**
```sql
SELECT 
    s.cnpj_cpf_socio,
    s.nome_socio,
    COUNT(DISTINCT CASE WHEN est.situacao_cadastral = '02' THEN s.cnpj END) as empresas_ativas,
    COUNT(DISTINCT CASE WHEN est.situacao_cadastral = '08' THEN s.cnpj END) as empresas_baixadas
FROM socios s
JOIN estabelecimento est ON s.cnpj = est.cnpj
GROUP BY s.cnpj_cpf_socio, s.nome_socio
HAVING empresas_baixadas > 0 AND empresas_ativas > 0
ORDER BY empresas_baixadas DESC, empresas_ativas DESC
```

**Resultado:** Pessoas com empresas ativas e baixadas

### 12. **Rede de 2º Grau (Sócios de Sócios)**
```sql
WITH socios_empresa AS (
    SELECT cnpj_cpf_socio, nome_socio
    FROM socios
    WHERE cnpj = ?
),
empresas_socios AS (
    SELECT DISTINCT s2.cnpj, s2.nome_socio as socio_2grau, s2.cnpj_cpf_socio as cpf_2grau
    FROM socios_empresa se
    JOIN socios s2 ON se.cnpj_cpf_socio = s2.cnpj_cpf_socio
    WHERE s2.cnpj != ?
)
SELECT 
    es.cnpj,
    e.razao_social,
    est.nome_fantasia,
    es.socio_2grau,
    es.cpf_2grau
FROM empresas_socios es
JOIN estabelecimento est ON es.cnpj = est.cnpj
JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
```

**Resultado:** Outras empresas dos sócios da empresa inicial

### 13. **Empresas por CNAE**
```sql
SELECT 
    est.cnae_fiscal,
    c.descricao as atividade,
    COUNT(*) as total_empresas,
    COUNT(DISTINCT est.municipio) as total_municipios
FROM estabelecimento est
LEFT JOIN cnae c ON est.cnae_fiscal = c.codigo
WHERE est.situacao_cadastral = '02'
GROUP BY est.cnae_fiscal, c.descricao
ORDER BY total_empresas DESC
LIMIT 100
```

**Resultado:** CNAEs mais comuns

### 14. **Concentração Geográfica**
```sql
SELECT 
    est.uf,
    m.descricao as municipio,
    COUNT(*) as total_empresas,
    COUNT(DISTINCT s.cnpj_cpf_socio) as total_socios
FROM estabelecimento est
LEFT JOIN municipio m ON est.municipio = m.codigo
LEFT JOIN socios s ON est.cnpj = s.cnpj
WHERE est.situacao_cadastral = '02'
GROUP BY est.uf, m.descricao
ORDER BY total_empresas DESC
LIMIT 100
```

**Resultado:** Concentração de empresas por cidade

### 15. **Timeline de Atividades**
```sql
SELECT 
    s.cnpj_cpf_socio,
    s.nome_socio,
    s.cnpj,
    e.razao_social,
    s.data_entrada_sociedade,
    est.data_inicio_atividades,
    est.situacao_cadastral,
    est.data_situacao_cadastral
FROM socios s
JOIN estabelecimento est ON s.cnpj = est.cnpj
JOIN empresas e ON est.cnpj_basico = e.cnpj_basico
WHERE s.cnpj_cpf_socio = ?
ORDER BY s.data_entrada_sociedade, est.data_inicio_atividades
```

**Resultado:** Linha do tempo de atividades empresariais de uma pessoa

## 📊 Estatísticas Disponíveis

### Volume de Dados
- ~50 milhões de empresas
- ~20 milhões de sócios únicos
- ~100 milhões de relacionamentos

### Campos Únicos para Triangulação
1. **CPF** (sem máscara)
2. **CNPJ** (14 dígitos)
3. **Nome completo**
4. **Endereço completo**
5. **Email**
6. **Telefone**
7. **CEP**
8. **Data de nascimento** (via faixa etária)
9. **Representante legal**

## 🚀 Implementação Necessária

### 1. Remover Censura
- Eliminar código que censura dados de empresários individuais
- Exibir TODOS os campos sem restrição

### 2. Criar Queries de Cruzamento
- Implementar todas as 15 queries acima
- Adicionar índices para performance
- Criar views materializadas

### 3. Interface CLI
- Menu de cruzamentos
- Exportação de resultados
- Visualização de redes

### 4. APIs REST
- Endpoints para cada tipo de cruzamento
- Paginação de resultados
- Filtros avançados
