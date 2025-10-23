# Guia de Instalação - RedeCNPJ Go

## Pré-requisitos

- Go 1.21 ou superior
- SQLite3
- Make (opcional, mas recomendado)

## Instalação do Go

### Linux (Ubuntu/Debian)
```bash
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### macOS
```bash
brew install go
```

### Windows
Baixe o instalador em: https://go.dev/dl/

## Instalação do Projeto

### 1. Clone o repositório (se ainda não fez)
```bash
git clone https://github.com/peder1981/rede-cnpj.git
cd rede-cnpj/RedeGO
```

### 2. Instale as dependências
```bash
make deps
# ou
go mod download
```

### 3. Configure os bancos de dados

Edite o arquivo `rede.ini` e configure os caminhos para os bancos de dados:

```ini
[BASE]
base_receita = bases/cnpj.db
base_rede = bases/rede.db
base_rede_search = bases/rede_search.db
```

### 4. Copie os bancos de dados

Copie os arquivos `.db` da pasta `../rede/bases/` para `./bases/`:

```bash
mkdir -p bases
cp ../rede/bases/cnpj_teste.db bases/
cp ../rede/bases/rede_teste.db bases/
```

### 5. Compile o projeto
```bash
make build
# ou
go build -o rede-cnpj ./cmd/server
```

### 6. Execute
```bash
make run
# ou
./rede-cnpj
```

O servidor estará disponível em: http://127.0.0.1:5000/rede/

## Opções de Linha de Comando

```bash
./rede-cnpj -h
```

Opções disponíveis:
- `-conf_file`: Arquivo de configuração (padrão: rede.ini)
- `-inicial`: CNPJ/CPF inicial para carregar
- `-camada`: Número de camadas (padrão: 1)
- `-porta_flask`: Porta do servidor (padrão: 5000)
- `-pasta`: Pasta de arquivos (padrão: arquivos)

Exemplo:
```bash
./rede-cnpj -conf_file=rede.ini -porta_flask=8080
```

## Compilação para Produção

Para compilar uma versão otimizada:

```bash
make build-prod
```

Para compilar para múltiplas plataformas:

```bash
make build-all
```

Isso criará binários em `./build/` para:
- Linux (amd64)
- Windows (amd64)
- macOS (amd64 e arm64)

## Estrutura de Diretórios

Após a instalação, a estrutura deve ser:

```
RedeGO/
├── bases/              # Bancos de dados SQLite
│   ├── cnpj.db
│   ├── rede.db
│   └── rede_search.db
├── arquivos/           # Arquivos do usuário
├── static/             # Arquivos estáticos (CSS, JS, imagens)
├── templates/          # Templates HTML
├── cmd/
│   └── server/         # Aplicação principal
├── internal/           # Código interno
├── pkg/                # Pacotes públicos
├── rede.ini            # Configuração
└── rede-cnpj           # Binário executável
```

## Troubleshooting

### Erro: "banco de dados não encontrado"
Verifique se os caminhos no `rede.ini` estão corretos e se os arquivos `.db` existem.

### Erro: "porta já em uso"
Altere a porta no `rede.ini` ou use a opção `-porta_flask`:
```bash
./rede-cnpj -porta_flask=8080
```

### Erro de compilação
Certifique-se de ter Go 1.21+ instalado:
```bash
go version
```

Limpe o cache e recompile:
```bash
make clean
make build
```

## Diferenças da Versão Python

### Vantagens
- **Performance**: 5-10x mais rápido em consultas
- **Memória**: Uso de memória 3-4x menor
- **Deploy**: Binário único, sem dependências externas
- **Concorrência**: Melhor handling de requisições simultâneas

### Limitações Atuais
- Algumas funcionalidades avançadas ainda não implementadas:
  - Busca com palavras-chave (spacy)
  - Exportação completa para i2 Chart Reader
  - Geração de mapas com Folium
  
Estas funcionalidades serão adicionadas em versões futuras.

## Desenvolvimento

Para contribuir com o desenvolvimento:

```bash
# Instalar ferramentas de desenvolvimento
make install-tools

# Executar testes
make test

# Verificar código
make vet
make lint

# Formatar código
make fmt
```

## Suporte

Para reportar bugs ou solicitar funcionalidades, abra uma issue em:
https://github.com/peder1981/rede-cnpj/issues
