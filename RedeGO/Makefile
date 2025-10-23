.PHONY: build run clean test install deps

# Variáveis
BINARY_NAME=rede-cnpj
MAIN_PATH=./cmd/server
BUILD_DIR=./build

# Build servidor de APIs
build:
	@echo "Compilando $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build concluído: $(BINARY_NAME)"

build-cli:
	@echo "Compilando rede-cnpj-cli..."
	@go build -o rede-cnpj-cli ./cmd/cli
	@echo "Build concluído: rede-cnpj-cli"

build-importer:
	@echo "Compilando rede-cnpj-importer..."
	@go build -o rede-cnpj-importer ./cmd/importer
	@echo "Build concluído: rede-cnpj-importer"

build-all-binaries: build build-cli build-importer
	@echo "Todos os binários compilados!"

# Build para produção (otimizado)
build-prod:
	@echo "Compilando $(BINARY_NAME) para produção..."
	@go build -ldflags="-s -w" -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build de produção concluído: $(BINARY_NAME)"

# Build para múltiplas plataformas
build-all:
	@echo "Compilando para múltiplas plataformas..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "Builds concluídos em $(BUILD_DIR)/"

# Executar
run: build
	@echo "Executando $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Executar com configuração customizada
run-custom: build
	@echo "Executando $(BINARY_NAME) com configuração customizada..."
	@./$(BINARY_NAME) -conf_file=$(CONFIG)

# Instalar dependências
deps:
	@echo "Baixando dependências..."
	@go mod download
	@go mod tidy
	@echo "Dependências instaladas"

# Atualizar dependências
update-deps:
	@echo "Atualizando dependências..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependências atualizadas"

# Testes
test:
	@echo "Executando testes..."
	@go test -v ./...

# Testes com cobertura
test-coverage:
	@echo "Executando testes com cobertura..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório de cobertura gerado: coverage.html"

# Linter
lint:
	@echo "Executando linter..."
	@golangci-lint run ./...

# Formatar código
fmt:
	@echo "Formatando código..."
	@go fmt ./...
	@gofmt -s -w .

# Verificar código
vet:
	@echo "Verificando código..."
	@go vet ./...

# Limpar arquivos gerados
clean:
	@echo "Limpando arquivos gerados..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Limpeza concluída"

# Instalar ferramentas de desenvolvimento
install-tools:
	@echo "Instalando ferramentas de desenvolvimento..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Ferramentas instaladas"

# Criar diretórios necessários
setup:
	@echo "Criando estrutura de diretórios..."
	@mkdir -p bases
	@mkdir -p arquivos
	@mkdir -p static
	@mkdir -p templates
	@touch bases/.gitkeep
	@touch arquivos/.gitkeep
	@echo "Estrutura criada"

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  make build         - Compila o projeto"
	@echo "  make build-prod    - Compila otimizado para produção"
	@echo "  make build-all     - Compila para múltiplas plataformas"
	@echo "  make run           - Compila e executa"
	@echo "  make deps          - Instala dependências"
	@echo "  make update-deps   - Atualiza dependências"
	@echo "  make test          - Executa testes"
	@echo "  make test-coverage - Executa testes com cobertura"
	@echo "  make lint          - Executa linter"
	@echo "  make fmt           - Formata código"
	@echo "  make vet           - Verifica código"
	@echo "  make clean         - Remove arquivos gerados"
	@echo "  make setup         - Cria estrutura de diretórios"
	@echo "  make help          - Exibe esta ajuda"
