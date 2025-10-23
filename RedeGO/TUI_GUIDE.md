# 🎮 RedeCNPJ TUI - Interface Estilo DOS

## Interface Completa com Menus Visuais

Interface de navegação interativa estilo DOS com menus suspensos, feedback visual claro e ajuda integrada.

## 🚀 Como Usar

```bash
make build-cli
./rede-cnpj-cli -conf_file=rede.ini
```

Digite o CNPJ inicial e navegue!

## 📺 Modos de Visualização

### 1. **MODO ÁRVORE** (Padrão)

Navegação hierárquica dos relacionamentos:

```
╔══════════════════════════════════════════════════════════════════════╗
║         🔍 RedeCNPJ - Navegação em Árvore                           ║
╚══════════════════════════════════════════════════════════════════════╝

→ ▶ 🏢 CRUISER INFORMATICA E SERVICOS LTDA
  ▶ 👤 VINICIUS D ANTONIO
  ▶ 👤 DENISE ALT PINTO

┌──────────────────────────────────────────────────────────────────────┐
│ ✓ Carregado: 4 nós, 3 ligações                                      │
└──────────────────────────────────────────────────────────────────────┘
┌──────────────────────────────────────────────────────────────────────┐
│ NAVEGAÇÃO: ↑↓ mover | → expandir | ← colapsar                       │
│ AÇÕES: [A]nalytics | [E]xportar | [F1/?] Ajuda | [ESC] Sair        │
└──────────────────────────────────────────────────────────────────────┘
```

**Comandos:**
- **↑/k** - Mover para cima
- **↓/j** - Mover para baixo
- **→/l/Enter** - Expandir nó (busca relacionamentos)
- **←/h** - Colapsar nó
- **a** - Modo Analytics
- **e** - Modo Exportar
- **F1/?** - Ajuda
- **ESC** - Sair

### 2. **MODO ANALYTICS**

Estatísticas completas do grafo:

```
╔══════════════════════════════════════════════════════════════════════╗
║         📊 RedeCNPJ - Estatísticas do Grafo                         ║
╚══════════════════════════════════════════════════════════════════════╝

┌─ RESUMO GERAL ─────────────────────────────────────────────────────┐
│ Total de Nós:                4                                     │
│ Total de Ligações:           3                                     │
│ Empresas (PJ):               1                                     │
│ Pessoas (PF):                3                                     │
│ Pessoas Externas (PE):       0                                     │
└────────────────────────────────────────────────────────────────────┘

┌─ MÉTRICAS DE REDE ─────────────────────────────────────────────────┐
│ Densidade:              0.2500                                     │
│ Grau Médio:             1.50                                       │
│ Componentes Conexos:         1                                     │
└────────────────────────────────────────────────────────────────────┘

┌─ NÓS MAIS CONECTADOS ──────────────────────────────────────────────┐
│  1. CRUISER INFORMATICA E SERVICOS LTDA      [  3 conexões]       │
│  2. VINICIUS D ANTONIO                       [  1 conexões]       │
└────────────────────────────────────────────────────────────────────┘

┌─ TIPOS DE RELACIONAMENTO ──────────────────────────────────────────┐
│ Sócio Administrador                         :      1               │
│ Sócio                                       :      2               │
└────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────┐
│ Pressione [Q] ou [BACKSPACE] para voltar                            │
└──────────────────────────────────────────────────────────────────────┘
```

**O que mostra:**
- Total de nós e ligações
- Distribuição por tipo (Empresas/Pessoas)
- Densidade da rede
- Grau médio
- Top 10 nós mais conectados
- Tipos de relacionamento

**Comandos:**
- **Q/Backspace** - Voltar ao modo árvore

### 3. **MODO EXPORTAR**

Menu visual para exportação:

```
╔══════════════════════════════════════════════════════════════════════╗
║         💾 RedeCNPJ - Exportar Dados                                ║
╚══════════════════════════════════════════════════════════════════════╝

Selecione o formato de exportação:

→ 📊 Excel (XLSX) - Arquivo completo com 3 planilhas
  📄 CSV - Nós (lista de entidades)
  📄 CSV - Arestas (lista de relacionamentos)
  📄 CSV - Estatísticas (resumo do grafo)

┌──────────────────────────────────────────────────────────────────────┐
│ ↑↓ Navegar | [ENTER] Exportar | [Q] Cancelar                        │
│ Arquivos serão salvos em: ./output/                                 │
└──────────────────────────────────────────────────────────────────────┘
```

**Formatos disponíveis:**

1. **Excel (XLSX)** - Arquivo completo com 3 planilhas:
   - Nós (entidades)
   - Arestas (relacionamentos)
   - Estatísticas

2. **CSV - Nós** - Lista de todas as entidades
3. **CSV - Arestas** - Lista de todos os relacionamentos
4. **CSV - Estatísticas** - Resumo estatístico

**Comandos:**
- **↑↓** - Navegar opções
- **Enter** - Exportar formato selecionado
- **Q** - Cancelar e voltar

**Arquivos salvos em:** `./output/`

### 4. **MODO AJUDA**

Ajuda completa integrada:

```
╔══════════════════════════════════════════════════════════════════════╗
║         ❓ RedeCNPJ - Ajuda                                          ║
╚══════════════════════════════════════════════════════════════════════╝

┌─ MODO ÁRVORE (Navegação) ──────────────────────────────────────────┐
│                                                                      │
│  ↑ / k          - Mover cursor para cima                            │
│  ↓ / j          - Mover cursor para baixo                           │
│  → / l / ENTER  - Expandir nó selecionado (busca relacionamentos)   │
│  ← / h          - Colapsar nó selecionado                           │
│  a              - Ver Analytics (estatísticas do grafo)             │
│  e              - Exportar dados (Excel ou CSV)                     │
│  F1 / ?         - Mostrar esta ajuda                                │
│  ESC            - Sair do programa                                  │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘

[... mais seções de ajuda ...]
```

**Comandos:**
- **F1/?** - Abrir ajuda
- **Q/Backspace/F1/?** - Fechar ajuda

## 🎯 Fluxo de Uso Típico

### 1. Iniciar e Navegar
```bash
./rede-cnpj-cli -conf_file=rede.ini
# Digite: 01212126000192
```

### 2. Expandir Relacionamentos
- Use **↓** para selecionar um nó
- Pressione **→** ou **Enter** para expandir
- Veja os relacionamentos aparecerem indentados

### 3. Ver Estatísticas
- Pressione **a** para Analytics
- Veja estatísticas completas
- Pressione **q** para voltar

### 4. Exportar Dados
- Pressione **e** para Exportar
- Use **↑↓** para escolher formato
- Pressione **Enter** para exportar
- Arquivo salvo em `./output/`

## 📊 Indicadores Visuais

### Ícones
- **🏢** - Empresa (Pessoa Jurídica)
- **👤** - Pessoa (Pessoa Física ou Estrangeira)
- **▶** - Nó pode ser expandido
- **▼** - Nó já está expandido
- **→** - Cursor (item selecionado)

### Mensagens de Status
- **✓** - Operação bem-sucedida
- **✗** - Erro
- **⏳** - Processando

### Exemplos:
```
✓ Carregado: 4 nós, 3 ligações
✓ Expandido: +10 nós
✓ Exportado: output/rede-cnpj.xlsx
✗ Carregue um grafo primeiro
⏳ Expandindo...
```

## 💡 Dicas

1. **Navegação Rápida**
   - Use **j/k** (estilo Vim) ao invés de setas
   - Use **h/l** para colapsar/expandir

2. **Antes de Exportar**
   - Expanda os nós que deseja incluir
   - Use Analytics para ver o tamanho do grafo
   - Escolha o formato adequado

3. **Performance**
   - Grafos grandes (>1000 nós) podem demorar
   - Use filtros se necessário
   - Exporte em partes se muito grande

4. **Arquivos de Saída**
   - Excel: Melhor para análise visual
   - CSV: Melhor para processamento
   - Verifique `./output/` após exportar

## 🔧 Troubleshooting

### "✗ Carregue um grafo primeiro"
- Você tentou Analytics/Exportar sem dados
- Solução: Expanda pelo menos o nó raiz primeiro

### "✗ Erro ao exportar"
- Verifique permissões da pasta `./output/`
- Solução: `mkdir -p output && chmod 755 output`

### Nada acontece ao expandir
- O nó pode não ter relacionamentos
- Aguarde a mensagem de status

### Interface cortada
- Terminal muito pequeno
- Solução: Redimensione para pelo menos 80x24

## ⌨️ Referência Rápida de Comandos

### Global (Qualquer Modo)
| Tecla | Ação |
|-------|------|
| **F1** ou **?** | Ajuda |
| **ESC** | Sair |

### Modo Árvore
| Tecla | Ação |
|-------|------|
| **↑** ou **k** | Mover para cima |
| **↓** ou **j** | Mover para baixo |
| **→** ou **l** ou **Enter** | Expandir nó |
| **←** ou **h** | Colapsar nó |
| **a** | Analytics |
| **e** | Exportar |

### Modo Analytics
| Tecla | Ação |
|-------|------|
| **q** ou **Backspace** | Voltar |

### Modo Exportar
| Tecla | Ação |
|-------|------|
| **↑↓** | Navegar opções |
| **Enter** | Exportar |
| **q** | Cancelar |

### Modo Ajuda
| Tecla | Ação |
|-------|------|
| **q** ou **Backspace** ou **F1** ou **?** | Fechar |

## 🎨 Características

- ✅ **Visual Clara** - Menus e bordas estilo DOS
- ✅ **Feedback Imediato** - Mensagens de status em tempo real
- ✅ **Ajuda Integrada** - F1 sempre disponível
- ✅ **Navegação Intuitiva** - Setas ou Vim (hjkl)
- ✅ **Exportação Funcional** - Arquivos salvos em ./output/
- ✅ **Analytics Detalhado** - Estatísticas completas
- ✅ **Sem Ambiguidade** - Sempre sabe onde está e o que fazer
