# 📊 Análise de Features do Projeto Python

## Funcionalidades Identificadas

### 1. **APIs REST** (rede.py)

#### Endpoints Principais
- ✅ `/rede/grafojson/<tipo>/<camada>/<cpfcnpj>` - Grafo JSON (JÁ IMPLEMENTADO)
- ✅ `/rede/dadosjson/<cpfcnpj>` - Dados detalhados (JÁ IMPLEMENTADO)
- ⚠️ `/rede/grafojson/links/<camada>/<numeroItens>/<valorMinimo>/<valorMaximo>/<cpfcnpj>` - Grafo com filtros de links
- ⚠️ `/rede/consulta_cnpj/` - Consulta CNPJ HTML
- ⚠️ `/rede/api/<tipo>/<cnpj>` - API genérica
- ⚠️ `/rede/api/caminhos` - API de caminhos entre entidades
- ⚠️ `/rede/mapa` - Geração de mapas geográficos
- ⚠️ `/rede/dadosemarquivo/<formato>` - Exportação (xlsx, pdf, etc)
- ⚠️ `/rede/busca_google` - Busca Google integrada
- ⚠️ `/rede/informacao/dados_publicos_cnpj_disponivel` - Info sobre dados disponíveis

#### Gestão de Arquivos
- `/rede/arquivos_json/<arquivopath>` - CRUD de arquivos JSON
- `/rede/arquivos_json_upload/<nomeArquivo>` - Upload de arquivos
- `/rede/arquivo_upload/` - Upload genérico
- `/rede/json_para_base/<comentario>` - Salvar JSON no banco
- `/rede/abrir_arquivo` - Abrir arquivo local

### 2. **Funcionalidades de Busca** (rede_sqlite_cnpj.py)

#### Busca por Nome
- `buscaPorNome(nome, limite)` - Busca FTS5 com:
  - Suporte a wildcards (* ?)
  - Match flexível (não precisa nome completo)
  - Busca em múltiplas tabelas
  - Limite configurável

#### Tipos de Busca
- Busca por razão social
- Busca por nome fantasia
- Busca por nome de sócio
- Busca por CNPJ/CPF
- Busca com curingas
- Busca aleatória (#TESTE#)

### 3. **Geração de Grafos** (rede_sqlite_cnpj.py)

#### Tipos de Grafo
- **rede** - Rede completa de relacionamentos
- **caminhos** - Caminhos entre entidades
- **caminhos-direto** - Apenas caminhos diretos
- **caminhos-comum** - Entidades em comum
- **links** - Baseado em tabela de links customizada

#### Parâmetros
- Camadas (profundidade)
- Limite de registros por camada
- Tempo máximo de consulta
- Filtros por tipo de ligação
- Inclusão de filiais

### 4. **Exportação de Dados**

#### Formatos Suportados
- **XLSX** - Excel
- **PDF** - Relatórios PDF
- **CSV** - Dados tabulares
- **JSON** - Dados estruturados
- **HTML** - Páginas web

#### Tipos de Exportação
- Dados de empresas
- Dados de sócios
- Grafos completos
- Relatórios customizados

### 5. **Integração com Busca Externa**

#### Google Search (rede_busca.py)
- Busca no Google
- Extração de palavras-chave
- NLP com Spacy (opcional)
- Análise de sites
- Filtro de sites a ignorar

#### DuckDuckGo (ddgs)
- Busca alternativa
- Sem rate limiting
- Mais privacidade

### 6. **Mapas Geográficos** (mapa.py)

#### Funcionalidades
- Geocoding de endereços
- Mapa de municípios
- Visualização geográfica
- Latitude/Longitude
- JSON com coordenadas

### 7. **Análise de Dados** (rede_sqlite_cnpj.py)

#### Dicionários
- Qualificação de sócios
- Motivos de situação cadastral
- CNAEs
- Natureza jurídica
- Situação cadastral
- Porte de empresa
- Países
- Municípios

#### Análises
- Contagem de vínculos
- Análise de caminhos
- Detecção de padrões
- Estatísticas de rede

### 8. **Features Especiais**

#### Flags de Análise
- `situacao_fiscal` - Situação fiscal
- `pep` - Pessoa Exposta Politicamente
- `ceis` - Cadastro de Empresas Inidôneas
- `cepim` - Cadastro de Entidades Privadas
- `cnep` - Cadastro Nacional de Empresas Punidas
- `acordo_leniência` - Acordos de leniência
- `ceaf` - Cadastro de Expulsões
- `pgfn-fgts` - Dívida FGTS
- `pgfn-sida` - Dívida ativa
- `pgfn-prev` - Dívida previdenciária
- `servidor_siape` - Servidores públicos

#### Rate Limiting
- Limite por minuto
- Limite por endpoint
- Proteção contra abuso

#### Segurança
- API Keys
- Verificação de usuário local
- Upload de arquivos seguro
- Validação de extensões

### 9. **Otimizações**

#### Cache
- `@lru_cache` em funções críticas
- Cache de dicionários
- Cache de consultas

#### Performance
- Consultas in-memory quando possível
- Índices otimizados
- Limite de tempo de consulta
- Limite de registros

#### Concorrência
- Lock para SQLite
- Suporte a uwsgi
- Threading lock
- Tabelas temporárias com prefixo

### 10. **Integração i2** (rede_i2.py)

#### i2 Analyst's Notebook
- Importação de dados
- Exportação para i2
- Formato ANB/ANX
- Análise de inteligência

## Priorização para Implementação em Go

### ✅ Já Implementado
1. Servidor APIs REST básico
2. Busca por CNPJ
3. Grafo JSON básico
4. Dados detalhados

### 🔥 Alta Prioridade
1. **Busca Avançada**
   - Busca por nome com wildcards
   - Busca FTS5 otimizada
   - Múltiplos critérios

2. **Tipos de Grafo**
   - Caminhos entre entidades
   - Filtros avançados
   - Grafos customizados

3. **Exportação**
   - Excel (XLSX)
   - CSV
   - JSON estruturado

4. **Análise de Dados**
   - Estatísticas de rede
   - Detecção de padrões
   - Relatórios automáticos

### 🟡 Média Prioridade
1. **Mapas Geográficos**
   - Geocoding
   - Visualização em mapa

2. **Flags de Análise**
   - PEP, CEIS, CNEP, etc
   - Integração com bases externas

3. **Gestão de Arquivos**
   - Upload/Download
   - Armazenamento de grafos

### 🔵 Baixa Prioridade
1. **Busca Externa**
   - Google/DuckDuckGo
   - NLP/Spacy

2. **Integração i2**
   - Formato ANB/ANX

3. **PDF**
   - Geração de relatórios PDF

## Tecnologias a Incorporar

### Bibliotecas Go Necessárias
- `github.com/xuri/excelize/v2` - Excel
- `github.com/jung-kurt/gofpdf` - PDF
- `github.com/paulmach/orb` - Geolocalização
- `github.com/go-echarts/go-echarts/v2` - Gráficos
- Rate limiting (já temos com Gin)

### Arquitetura Proposta
```
internal/
├── analytics/      # Análises e estatísticas
├── export/         # Exportação (xlsx, csv, pdf)
├── geocoding/      # Mapas e geolocalização
├── graph/          # Tipos de grafo avançados
├── search/         # Busca avançada
└── flags/          # Flags de análise (PEP, CEIS, etc)
```

## Próximos Passos

1. Implementar busca avançada com wildcards
2. Adicionar tipos de grafo (caminhos, filtros)
3. Implementar exportação Excel/CSV
4. Criar módulo de analytics
5. Adicionar flags de análise
6. Implementar geocoding básico
