# Índice de Documentação - RedeCNPJ Go

Guia completo de navegação pela documentação do projeto.

## 🚀 Para Começar

1. **[QUICKSTART.md](QUICKSTART.md)** - Comece aqui! 
   - Instalação em 3 passos
   - Primeiros testes
   - Exemplos básicos
   - ⏱️ Tempo de leitura: 5 minutos

2. **[README.md](README.md)** - Visão Geral
   - Descrição do projeto
   - Estrutura de diretórios
   - Requisitos
   - Diferenças da versão Python
   - ⏱️ Tempo de leitura: 3 minutos

## 📖 Documentação Técnica

3. **[INSTALL.md](INSTALL.md)** - Instalação Detalhada
   - Instalação do Go
   - Configuração completa
   - Opções de linha de comando
   - Compilação para produção
   - Troubleshooting
   - ⏱️ Tempo de leitura: 10 minutos

4. **[ARCHITECTURE.md](ARCHITECTURE.md)** - Arquitetura
   - Estrutura do projeto
   - Camadas da aplicação
   - Padrões de design
   - Fluxo de dados
   - Concorrência
   - Performance
   - Segurança
   - Extensibilidade
   - ⏱️ Tempo de leitura: 20 minutos

## 🔄 Migração

5. **[MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)** - Guia de Migração
   - Comparação Python vs Go
   - Diferenças de implementação
   - Migração passo a passo
   - Compatibilidade de API
   - Scripts de migração
   - ⏱️ Tempo de leitura: 15 minutos

## 📊 Resumo

6. **[SUMMARY.md](SUMMARY.md)** - Resumo Completo
   - Estatísticas do projeto
   - Funcionalidades implementadas
   - Melhorias de performance
   - Tecnologias utilizadas
   - Roadmap
   - ⏱️ Tempo de leitura: 10 minutos

## 📁 Estrutura de Arquivos

### Código Fonte

```
RedeGO/
├── cmd/
│   └── server/
│       └── main.go                    # Ponto de entrada
│
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuração
│   ├── database/
│   │   └── database.go                # Acesso a dados
│   ├── handlers/
│   │   └── handlers.go                # Handlers HTTP
│   ├── models/
│   │   └── models.go                  # Modelos de dados
│   ├── services/
│   │   ├── rede_service.go            # Lógica de negócio
│   │   └── export_service.go          # Exportação
│   └── utils/
│       └── utils.go                   # Utilitários
│
└── pkg/
    └── cpfcnpj/
        ├── validator.go               # Validação CPF/CNPJ
        └── validator_test.go          # Testes
```

### Documentação

```
RedeGO/
├── INDEX.md                           # Este arquivo
├── QUICKSTART.md                      # Início rápido
├── README.md                          # Visão geral
├── INSTALL.md                         # Instalação
├── ARCHITECTURE.md                    # Arquitetura
├── MIGRATION_GUIDE.md                 # Migração
├── SUMMARY.md                         # Resumo
└── SCRIPTS.md                         # Scripts shell
```

### Scripts

```
RedeGO/
├── start.sh                           # Iniciar aplicação
├── stop.sh                            # Parar aplicação
└── restart.sh                         # Reiniciar aplicação
```

### Configuração

```
RedeGO/
├── go.mod                             # Dependências Go
├── rede.ini                           # Configuração app
├── Makefile                           # Automação
└── .gitignore                         # Git ignore
```

## 🎯 Guias por Objetivo

### Quero instalar rapidamente
→ [QUICKSTART.md](QUICKSTART.md)

### Quero entender a arquitetura
→ [ARCHITECTURE.md](ARCHITECTURE.md)

### Quero migrar de Python
→ [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)

### Quero ver estatísticas e performance
→ [SUMMARY.md](SUMMARY.md)

### Quero instalar em produção
→ [INSTALL.md](INSTALL.md) (seção "Compilação para Produção")

### Quero contribuir
→ [ARCHITECTURE.md](ARCHITECTURE.md) (seção "Contribuindo")

## 📚 Guias por Nível

### 👶 Iniciante
1. [QUICKSTART.md](QUICKSTART.md)
2. [README.md](README.md)
3. [INSTALL.md](INSTALL.md)

### 🧑 Intermediário
1. [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)
2. [SUMMARY.md](SUMMARY.md)
3. [ARCHITECTURE.md](ARCHITECTURE.md) (seções básicas)

### 👨‍💻 Avançado
1. [ARCHITECTURE.md](ARCHITECTURE.md) (completo)
2. Código fonte em `internal/`
3. Testes em `pkg/cpfcnpj/validator_test.go`

## 🔍 Busca Rápida

### Instalação
- Pré-requisitos: [INSTALL.md](INSTALL.md#pré-requisitos)
- Instalação do Go: [INSTALL.md](INSTALL.md#instalação-do-go)
- Compilação: [INSTALL.md](INSTALL.md#instalação-do-projeto)

### Configuração
- Arquivo rede.ini: [INSTALL.md](INSTALL.md#configure-os-bancos-de-dados)
- Linha de comando: [INSTALL.md](INSTALL.md#opções-de-linha-de-comando)
- Variáveis: [ARCHITECTURE.md](ARCHITECTURE.md#configuração)

### Uso
- Exemplos básicos: [QUICKSTART.md](QUICKSTART.md#exemplos-de-uso)
- API endpoints: [ARCHITECTURE.md](ARCHITECTURE.md#principais-endpoints)
- Exportação: [SUMMARY.md](SUMMARY.md#exportação-90)

### Performance
- Benchmarks: [SUMMARY.md](SUMMARY.md#melhorias-de-performance)
- Otimizações: [ARCHITECTURE.md](ARCHITECTURE.md#otimizações-implementadas)
- Comparação: [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md#performance)

### Desenvolvimento
- Estrutura: [ARCHITECTURE.md](ARCHITECTURE.md#estrutura-do-projeto)
- Padrões: [ARCHITECTURE.md](ARCHITECTURE.md#padrões-de-design)
- Testes: [SUMMARY.md](SUMMARY.md#testes)

## 🛠️ Comandos Rápidos

### Instalação
```bash
make deps          # Instalar dependências
make build         # Compilar
make run           # Executar
```

### Desenvolvimento
```bash
make test          # Testes
make fmt           # Formatar
make vet           # Verificar
make clean         # Limpar
```

### Produção
```bash
make build-prod    # Build otimizado
make build-all     # Múltiplas plataformas
```

Ver [Makefile](Makefile) para todos os comandos.

## 📞 Suporte

### Problemas Comuns
- Instalação: [INSTALL.md](INSTALL.md#troubleshooting)
- Migração: [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md#troubleshooting)
- Uso: [QUICKSTART.md](QUICKSTART.md#troubleshooting)

### Reportar Bugs
- GitHub Issues: https://github.com/peder1981/rede-cnpj/issues
- Incluir: versão Go, SO, logs, passos para reproduzir

### Contribuir
- Fork o repositório
- Leia [ARCHITECTURE.md](ARCHITECTURE.md#contribuindo)
- Envie Pull Request

## 📊 Estatísticas

- **Arquivos de código**: 10 (.go)
- **Linhas de código**: ~2,000
- **Arquivos de doc**: 7 (.md)
- **Testes**: 1 arquivo
- **Cobertura**: 60-95% (por módulo)

## 🗺️ Roadmap

Ver [SUMMARY.md](SUMMARY.md#roadmap) para planos futuros.

## ✅ Checklist de Leitura

Para dominar o projeto, leia nesta ordem:

- [ ] QUICKSTART.md - Instalação básica
- [ ] README.md - Visão geral
- [ ] INSTALL.md - Instalação completa
- [ ] SUMMARY.md - Estatísticas e features
- [ ] MIGRATION_GUIDE.md - Diferenças Python/Go
- [ ] ARCHITECTURE.md - Arquitetura detalhada
- [ ] Código fonte - Implementação

## 🎓 Recursos Adicionais

### Documentação Externa
- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [SQLite Go Driver](https://github.com/mattn/go-sqlite3)

### Projeto Original
- [RedeCNPJ Python](https://github.com/rictom/rede-cnpj)
- [Dados Abertos CNPJ](https://dados.gov.br/dados/conjuntos-dados/cadastro-nacional-da-pessoa-juridica---cnpj)

## 📝 Convenções

### Documentação
- **Negrito**: Conceitos importantes
- `código`: Comandos e código
- → Navegação/referência
- ✅ Implementado
- ⚠️ Parcial
- ❌ Não implementado
- ⏱️ Tempo estimado

### Código
- Seguir [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Usar `gofmt` para formatação
- Comentários em português no código de negócio
- Comentários em inglês em código genérico

## 🏆 Créditos

- **Projeto Original**: [rictom](https://github.com/rictom)
- **Conversão Go**: Desenvolvido para melhor performance e manutenibilidade
- **Dados**: Receita Federal do Brasil

---

**Última atualização**: 2025-01-21

**Versão**: 1.0.0

**Licença**: Ver [LICENSE](../LICENSE)
