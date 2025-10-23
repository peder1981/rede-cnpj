# 📥 RedeCNPJ Importer - Guia Completo

## Sistema de Importação de Dados da Receita Federal

Importador completo e eficiente dos dados públicos de CNPJ da Receita Federal, escrito em Go.

## 🎯 Funcionalidades

1. **Download Automático** - Baixa arquivos ZIP da Receita Federal
2. **Processamento Paralelo** - Descompacta e processa CSVs em paralelo
3. **Criação de Banco SQLite** - Gera `cnpj.db` otimizado
4. **Tabelas de Ligação** - Cria `rede.db` com relacionamentos
5. **Índices de Busca** - Cria `rede_search.db` com FTS5

## 🚀 Como Usar

### Compilar

```bash
make build-importer
```

### Executar Processo Completo

```bash
./rede-cnpj-importer -all
```

Este comando executa:
1. Download dos arquivos ZIP (~37 arquivos, ~15GB)
2. Processamento e criação do `cnpj.db` (~50GB)
3. Criação de tabelas de ligação `rede.db` (~20GB)
4. Criação de índices de busca `rede_search.db` (~5GB)

**Tempo estimado:** 2-4 horas (depende da conexão e hardware)

### Executar Etapas Individuais

#### 1. Apenas Download

```bash
./rede-cnpj-importer -download
```

Baixa os arquivos ZIP para `dados-publicos-zip/`

#### 2. Apenas Processamento

```bash
./rede-cnpj-importer -process
```

Processa arquivos já baixados e cria `cnpj.db`

#### 3. Apenas Tabelas de Ligação

```bash
./rede-cnpj-importer -links
```

Cria `rede.db` a partir do `cnpj.db`

#### 4. Apenas Índices de Busca

```bash
./rede-cnpj-importer -search
```

Cria `rede_search.db` com índices FTS5

## 📁 Estrutura de Diretórios

```
RedeGO/
├── dados-publicos-zip/     # Arquivos ZIP baixados
├── dados-publicos/         # CSVs descompactados (temporário)
└── bases/                  # Bancos de dados gerados
    ├── cnpj.db            # Base completa de CNPJ
    ├── rede.db            # Tabelas de ligação
    └── rede_search.db     # Índices de busca
```

## 📊 Bancos de Dados Gerados

### 1. cnpj.db (~50GB)

Tabelas principais:
- `empresas` - Dados das empresas (matrizes)
- `estabelecimento` - Estabelecimentos (matrizes + filiais)
- `socios` - Sócios e representantes
- `simples` - Optantes do Simples Nacional
- `cnae` - Códigos CNAE
- `municipio` - Códigos de municípios
- `natureza_juridica` - Naturezas jurídicas
- `pais` - Códigos de países
- `qualificacao_socio` - Qualificações de sócios
- `motivo` - Motivos de situação cadastral

### 2. rede.db (~20GB)

Tabela principal:
- `ligacao` - Relacionamentos entre entidades
  - `id1` - Origem (PJ_, PF_, PE_)
  - `id2` - Destino (PJ_, PF_, PE_)
  - `descricao` - Tipo de relacionamento
  - `comentario` - Base de origem

Tipos de ligação:
- **PJ → PJ** - Empresa sócia de empresa
- **PF → PJ** - Pessoa física sócia de empresa
- **PE → PJ** - Empresa estrangeira sócia
- **PF → PE** - Representante legal de empresa estrangeira
- **PJ → PJ** - Filial → Matriz

### 3. rede_search.db (~5GB)

Tabela FTS5:
- `id_search` - Índice full-text para busca rápida
  - Razões sociais
  - Nomes fantasia
  - Nomes de sócios
  - CNPJs/CPFs

## ⚙️ Requisitos de Sistema

### Mínimo
- **RAM:** 8GB
- **Disco:** 100GB livres
- **CPU:** 2 cores
- **Conexão:** 10 Mbps

### Recomendado
- **RAM:** 16GB+
- **Disco:** 200GB livres (SSD)
- **CPU:** 4+ cores
- **Conexão:** 50+ Mbps

## 🔄 Fluxo de Processamento

```
1. DOWNLOAD
   ├─ Busca última referência no site da Receita
   ├─ Lista arquivos ZIP disponíveis
   ├─ Download paralelo (5 simultâneos)
   └─ Salva em dados-publicos-zip/

2. PROCESSAMENTO
   ├─ Descompacta ZIPs
   ├─ Lê CSVs (encoding latin1, separador ;)
   ├─ Importa para SQLite
   │  ├─ empresas (7 colunas)
   │  ├─ estabelecimento (31 colunas)
   │  ├─ socios (12 colunas)
   │  └─ simples (7 colunas)
   ├─ Cria índices
   └─ Gera cnpj.db

3. LIGAÇÃO
   ├─ Anexa cnpj.db
   ├─ Cria tabela ligacao1 (temporária)
   │  ├─ PJ → PJ (sócios PJ)
   │  ├─ PF → PJ (sócios PF)
   │  ├─ PE → PJ (sócios estrangeiros)
   │  ├─ PF → PE (representantes)
   │  └─ PJ → PJ (filiais)
   ├─ Remove duplicatas
   ├─ Cria índices
   └─ Gera rede.db

4. INDEXAÇÃO
   ├─ Anexa cnpj.db e rede.db
   ├─ Cria tabela FTS5
   ├─ Indexa razões sociais
   ├─ Indexa nomes fantasia
   ├─ Indexa nomes de sócios
   └─ Gera rede_search.db
```

## 📈 Performance

### Download
- **Paralelo:** 5 downloads simultâneos
- **Tempo:** 30-60 min (depende da conexão)
- **Tamanho:** ~15GB

### Processamento
- **Batch:** 100.000 registros por transação
- **Tempo:** 1-2 horas
- **Registros:** ~50 milhões

### Ligação
- **Tempo:** 30-60 min
- **Ligações:** ~100 milhões

### Indexação
- **FTS5:** Full-text search otimizado
- **Tempo:** 20-40 min
- **Entradas:** ~60 milhões

## 🛠️ Otimizações Implementadas

1. **SQLite WAL Mode** - Write-Ahead Logging
2. **Batch Inserts** - Transações em lote
3. **Prepared Statements** - Queries pré-compiladas
4. **Parallel Downloads** - Downloads simultâneos
5. **Memory Cache** - Cache de 64MB
6. **Lazy Quotes** - Parsing CSV tolerante

## 🔍 Verificação

Após a importação, verifique:

```bash
# Tamanho dos bancos
ls -lh bases/*.db

# Quantidade de registros
sqlite3 bases/cnpj.db "SELECT COUNT(*) FROM empresas"
sqlite3 bases/rede.db "SELECT COUNT(*) FROM ligacao"
sqlite3 bases/rede_search.db "SELECT COUNT(*) FROM id_search"
```

## 📝 Logs

O importador exibe progresso em tempo real:

```
[1/37] ⬇️  Baixando Empresas0.zip...
[1/37] ✅ Empresas0.zip concluído
    Importando K3241.K03200Y0.D40114.EMPRECSV para tabela empresas...
      100000 registros...
      200000 registros...
      ✅ 5847852 registros importados
```

## ⚠️ Notas Importantes

1. **Espaço em Disco:** Certifique-se de ter pelo menos 100GB livres
2. **Conexão:** Download pode levar horas em conexões lentas
3. **Interrupção:** Se interrompido, delete os bancos e recomece
4. **Atualização:** Execute mensalmente para dados atualizados
5. **Backup:** Faça backup dos bancos após importação

## 🔄 Atualização Mensal

```bash
# 1. Remove bancos antigos
rm -f bases/cnpj.db bases/rede.db bases/rede_search.db

# 2. Remove arquivos antigos
rm -rf dados-publicos-zip/* dados-publicos/*

# 3. Executa importação completa
./rede-cnpj-importer -all
```

## 🆘 Troubleshooting

### Erro: "disk I/O error"
- **Causa:** Disco cheio
- **Solução:** Libere espaço (mínimo 100GB)

### Erro: "database is locked"
- **Causa:** Outro processo usando o banco
- **Solução:** Feche outros programas e tente novamente

### Download lento
- **Causa:** Conexão lenta ou site da Receita sobrecarregado
- **Solução:** Tente em horário alternativo

### Erro de memória
- **Causa:** RAM insuficiente
- **Solução:** Feche outros programas ou aumente swap

## 📚 Referências

- **Dados Públicos CNPJ:** https://www.gov.br/receitafederal/pt-br/assuntos/orientacao-tributaria/cadastros/consultas/dados-publicos-cnpj
- **Layout dos Arquivos:** Disponível no site da Receita Federal
- **SQLite FTS5:** https://www.sqlite.org/fts5.html
