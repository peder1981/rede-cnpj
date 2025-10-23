# 📁 Estrutura do Projeto RedeCNPJ Go

## Organização de Diretórios

```
RedeGO/
├── bin/                          # Binários compilados (gitignored)
│   ├── rede-cnpj                 # Servidor de APIs
│   ├── rede-cnpj-cli             # Interface TUI
│   └── rede-cnpj-importer        # Importador de dados
│
├── cmd/                          # Aplicações principais
│   ├── server/                   # Servidor de APIs REST
│   │   └── main.go
│   ├── cli/                      # Interface TUI
│   │   ├── main.go
│   │   ├── tui.go
│   │   └── tui_crossdata.go
│   └── importer/                 # Importador de dados
│       └── main.go
│
├── internal/                     # Código interno da aplicação
│   ├── analytics/                # Análise de grafos e estatísticas
│   │   └── stats.go
│   ├── config/                   # Configuração
│   │   └── config.go
│   ├── crossdata/                # Cruzamento de dados (SEM CENSURA)
│   │   └── queries.go
│   ├── database/                 # Acesso a banco de dados
│   │   └── database.go
│   ├── export/                   # Exportação (Excel, CSV)
│   │   ├── excel.go
│   │   └── csv.go
│   ├── forensics/                # Ferramentas forenses
│   │   └── investigator.go
│   ├── graph/                    # Tipos de grafo avançados
│   │   └── types.go
│   ├── handlers/                 # Handlers HTTP
│   │   ├── analytics.go
│   │   ├── crossdata.go
│   │   ├── export.go
│   │   ├── forensics.go
│   │   ├── graph.go
│   │   ├── handler.go
│   │   └── search.go
│   ├── importer/                 # Lógica de importação
│   │   ├── importer.go
│   │   ├── downloader.go
│   │   ├── processor.go
│   │   ├── linker.go
│   │   └── indexer.go
│   ├── models/                   # Modelos de dados
│   │   └── models.go
│   ├── search/                   # Busca avançada (FTS5)
│   │   └── advanced.go
│   ├── services/                 # Lógica de negócio
│   │   └── rede_service.go
│   └── utils/                    # Utilitários
│       └── utils.go
│
├── docs/                         # Documentação
│   ├── api/                      # Documentação de APIs
│   │   ├── API_COMPLETE.md       # Todas as APIs
│   │   ├── CROSSDATA_API.md      # APIs de cruzamento
│   │   └── FORENSICS_TOOLKIT.md  # Ferramentas forenses
│   └── guides/                   # Guias de uso
│       ├── TUI_GUIDE.md          # Guia da TUI
│       ├── IMPORTER_GUIDE.md     # Guia do importador
│       ├── FEATURES_ANALYSIS.md  # Análise de features
│       ├── CROSSDATA_SUMMARY.md  # Resumo de cruzamentos
│       └── DATABASE_ANALYSIS.md  # Análise dos bancos
│
├── doc/                          # Documentação técnica original
│   ├── QUICKSTART.md
│   ├── INDEX.md
│   ├── INSTALL.md
│   ├── ARCHITECTURE.md
│   ├── MIGRATION_GUIDE.md
│   └── SUMMARY.md
│
├── scripts/                      # Scripts de automação
│   ├── start.sh                  # Iniciar servidor
│   ├── stop.sh                   # Parar servidor
│   └── restart.sh                # Reiniciar servidor
│
├── bases/                        # Bancos de dados SQLite
│   ├── cnpj.db                   # Dados da Receita Federal
│   ├── rede.db                   # Rede de relacionamentos
│   └── rede_search.db            # Índice FTS5
│
├── arquivos/                     # Arquivos JSON salvos
├── output/                       # Arquivos exportados
├── logs/                         # Logs da aplicação
├── static/                       # Arquivos estáticos
├── templates/                    # Templates HTML
│
├── go.mod                        # Dependências Go
├── go.sum                        # Checksums de dependências
├── Makefile                      # Automação de build
├── .gitignore                    # Arquivos ignorados pelo Git
├── .env.example                  # Exemplo de variáveis de ambiente
├── rede.ini                      # Configuração principal
├── README.md                     # Documentação principal
└── PROJECT_STRUCTURE.md          # Este arquivo
```

## 📊 Estatísticas do Projeto

### Código Fonte
- **Go:** ~8.000 linhas
- **Documentação:** ~5.000 linhas
- **Total:** ~13.000 linhas

### Módulos Principais
1. **Server** - Servidor de APIs REST (Gin)
2. **CLI** - Interface TUI (Bubble Tea)
3. **Importer** - Importador de dados
4. **Analytics** - Análise de grafos
5. **CrossData** - Cruzamento de dados
6. **Forensics** - Ferramentas forenses
7. **Export** - Exportação (Excel/CSV)
8. **Search** - Busca avançada (FTS5)

### APIs Implementadas
- **15** endpoints básicos
- **12** endpoints de cruzamento
- **6** endpoints forenses
- **Total:** 33 endpoints REST

## 🔧 Comandos Principais

### Build
```bash
make build              # Compila servidor
make build-cli          # Compila CLI
make build-importer     # Compila importador
make build-all-binaries # Compila tudo
```

### Execução
```bash
./bin/rede-cnpj                 # Servidor
./bin/rede-cnpj-cli             # CLI
./bin/rede-cnpj-importer        # Importador
```

### Scripts
```bash
./scripts/start.sh              # Inicia servidor
./scripts/stop.sh               # Para servidor
./scripts/restart.sh            # Reinicia servidor
```

### Limpeza
```bash
make clean                      # Remove binários
rm -rf bin/ output/ logs/       # Limpeza completa
```

## 📦 Dependências Principais

### Runtime
- `github.com/gin-gonic/gin` - Framework web
- `github.com/mattn/go-sqlite3` - Driver SQLite
- `github.com/charmbracelet/bubbletea` - Framework TUI
- `github.com/xuri/excelize/v2` - Exportação Excel

### Build
- Go 1.21+
- Make (opcional)
- Git

## 🗄️ Bancos de Dados

### cnpj.db (~200GB)
- 50M+ empresas
- 20M+ sócios
- Dados da Receita Federal

### rede.db (~10GB)
- 100M+ relacionamentos
- Grafo de conexões

### rede_search.db (~5GB)
- Índice FTS5
- Busca full-text

## 📝 Convenções

### Nomenclatura
- **Arquivos:** snake_case.go
- **Pacotes:** lowercase
- **Tipos:** PascalCase
- **Funções:** camelCase (privadas) / PascalCase (públicas)
- **Constantes:** UPPER_CASE

### Estrutura de Código
```go
package nome

import (...)

// Tipos
type Struct struct {...}

// Constantes
const (...)

// Variáveis
var (...)

// Funções públicas
func PublicFunc() {...}

// Funções privadas
func privateFunc() {...}
```

## 🚀 Fluxo de Desenvolvimento

1. **Modificar código** em `internal/` ou `cmd/`
2. **Compilar** com `make build-all-binaries`
3. **Testar** executando binários em `bin/`
4. **Documentar** em `docs/`
5. **Commitar** mudanças

## 📊 Tamanho dos Módulos

| Módulo | Linhas | Arquivos |
|--------|--------|----------|
| forensics | 600+ | 1 |
| crossdata | 500+ | 1 |
| importer | 400+ | 5 |
| analytics | 300+ | 1 |
| handlers | 800+ | 6 |
| services | 200+ | 1 |
| cli | 700+ | 3 |
| **Total** | **~8.000** | **~30** |

## 🎯 Próximos Passos

### Organização
- ✅ Binários em `bin/`
- ✅ Documentação em `docs/`
- ✅ Scripts em `scripts/`
- ✅ Estrutura limpa e organizada

### Desenvolvimento
- ⏳ Testes unitários
- ⏳ CI/CD
- ⏳ Docker
- ⏳ Kubernetes

### Features
- ⏳ Dashboard web
- ⏳ Machine Learning
- ⏳ Grafos interativos
- ⏳ Alertas automáticos
