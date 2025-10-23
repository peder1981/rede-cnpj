# 🎮 RedeCNPJ TUI - Interface Interativa

## Navegação por Árvore com Setas

Interface de navegação interativa onde você pode **expandir cada nó individualmente** usando as setas do teclado.

## 🚀 Como Usar

### 1. Compilar

```bash
make build-cli
```

### 2. Executar

```bash
./rede-cnpj-cli -conf_file=rede.ini
```

### 3. Digite o CNPJ/CPF inicial

```
🔍 Digite o CNPJ/CPF inicial: 01212126000192
```

### 4. Navegue pela Árvore

```
╔════════════════════════════════════════════════════════════════╗
║         🔍 RedeCNPJ - Navegação Interativa                    ║
╚════════════════════════════════════════════════════════════════╝

📊 Use ↑↓ para navegar, →/Enter para expandir, ← para colapsar

→ ▶ 🏢 CRUISER INFORMATICA E SERVICOS LTDA
  ▶ 👤 VINICIUS D ANTONIO
  ▶ 👤 DENISE ALT PINTO
  ▶ 👤 RAFAEL ALT D ANTONIO

───────────────────────────────────────────────────────────────
Comandos: ↑↓ navegar | →/Enter expandir | ← colapsar | q sair
```

## 🎯 Controles

### Navegação
- **↑** ou **k** - Move para cima
- **↓** ou **j** - Move para baixo

### Expansão
- **→** ou **l** ou **Enter** ou **Espaço** - Expande o nó selecionado
- **←** ou **h** - Colapsa o nó selecionado

### Sair
- **q** ou **Ctrl+C** - Sair do programa

## 📖 Como Funciona

### 1. Nó Raiz (Inicial)
```
→ ▶ 🏢 CRUISER INFORMATICA E SERVICOS LTDA
```
- **→** = Cursor (nó selecionado)
- **▶** = Pode ser expandido
- **🏢** = Empresa (PJ)

### 2. Expandir Nó (pressione →)
```
→ ▼ 🏢 CRUISER INFORMATICA E SERVICOS LTDA
    ▶ 👤 VINICIUS D ANTONIO (Sócio Administrador)
    ▶ 👤 DENISE ALT PINTO (Sócio)
    ▶ 👤 RAFAEL ALT D ANTONIO (Sócio)
```
- **▼** = Nó expandido
- Filhos aparecem indentados

### 3. Expandir Qualquer Filho
```
  ▼ 🏢 CRUISER INFORMATICA E SERVICOS LTDA
    → ▼ 👤 VINICIUS D ANTONIO
        ▶ 🏢 OUTRA EMPRESA LTDA
        ▶ 🏢 MAIS UMA EMPRESA SA
    ▶ 👤 DENISE ALT PINTO
```
- Cada nó pode ser expandido independentemente
- Navegação em árvore completa

### 4. Colapsar (pressione ←)
```
→ ▶ 🏢 CRUISER INFORMATICA E SERVICOS LTDA
```
- Volta ao estado inicial
- Filhos são ocultados

## ✨ Características

- ✅ **Navegação Intuitiva** - Setas do teclado
- ✅ **Expansão Individual** - Cada nó é independente
- ✅ **Árvore Infinita** - Expanda quantos níveis quiser
- ✅ **Visual Limpo** - Indentação clara
- ✅ **Ícones** - 🏢 Empresas, 👤 Pessoas
- ✅ **Leve** - Interface em modo texto

## 🔄 Fluxo de Navegação

```
1. Inicia com nó raiz
   ▶ 🏢 EMPRESA

2. Expande raiz (→)
   ▼ 🏢 EMPRESA
     ▶ 👤 PESSOA A
     ▶ 👤 PESSOA B

3. Navega para PESSOA A (↓)
   ▼ 🏢 EMPRESA
   → ▶ 👤 PESSOA A
     ▶ 👤 PESSOA B

4. Expande PESSOA A (→)
   ▼ 🏢 EMPRESA
   → ▼ 👤 PESSOA A
       ▶ 🏢 EMPRESA 2
       ▶ 🏢 EMPRESA 3
     ▶ 👤 PESSOA B

5. Navega para EMPRESA 2 (↓)
   ▼ 🏢 EMPRESA
     ▼ 👤 PESSOA A
     → ▶ 🏢 EMPRESA 2
       ▶ 🏢 EMPRESA 3
     ▶ 👤 PESSOA B

6. Expande EMPRESA 2 (→)
   E assim por diante...
```

## 💡 Dicas

1. **Explore livremente** - Cada nó pode ser uma nova raiz
2. **Sem limites** - Expanda quantos níveis precisar
3. **Colapsar** - Use ← para limpar a visualização
4. **Navegação rápida** - Use j/k (estilo Vim) se preferir

## 🎨 Legenda

- **→** - Cursor (nó selecionado)
- **▶** - Nó pode ser expandido
- **▼** - Nó está expandido
- **🏢** - Empresa (Pessoa Jurídica)
- **👤** - Pessoa (Pessoa Física)
- **Indentação** - Nível na árvore
