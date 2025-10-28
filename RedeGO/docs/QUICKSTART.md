# Quick Start - RedeCNPJ Go

Guia rápido para começar a usar o RedeCNPJ em Go em 5 minutos.

## 📋 Pré-requisitos

- Go 1.21+ instalado
- Git
- 100MB de espaço livre

## 🚀 Instalação em 3 Passos

### 1. Instalar Dependências

```bash
cd /media/peder/DATA/rede-cnpj/RedeGO
make deps
```

Ou manualmente:
```bash
go mod download
go mod tidy
```

### 2. Copiar Bancos de Dados de Teste

```bash
# Criar diretório
mkdir -p bases

# Copiar bancos de teste
cp ../rede/bases/cnpj_teste.db bases/
cp ../rede/bases/rede_teste.db bases/
```

### 3. Compilar e Executar

```bash
make run
```

Ou manualmente:
```bash
go build -o rede-cnpj ./cmd/server
./rede-cnpj
```

## 🌐 Acessar a Aplicação

Abra seu navegador em:
```
http://127.0.0.1:5000/rede/
```

## ✅ Verificar Instalação

### Teste 1: Status da API
```bash
curl http://127.0.0.1:5000/rede/api/status
```

Resposta esperada:
```json
{
  "status": "ok",
  "version": "1.0.0",
  "message": "RedeCNPJ API em Go"
}
```

### Teste 2: Buscar Dados de CNPJ
```bash
curl http://127.0.0.1:5000/rede/dadosjson/12345678000190
```

### Teste 3: Buscar por Nome
```bash
curl "http://127.0.0.1:5000/rede/busca?q=EMPRESA&limite=5"
```

## 🎯 Exemplos de Uso

### Buscar Rede de Relacionamentos

```bash
curl -X POST http://127.0.0.1:5000/rede/grafojson/cnpj/2/12345678000190 \
  -H "Content-Type: application/json" \
  -d '["12345678000190"]'
```

### Exportar para Excel

```bash
curl -X POST http://127.0.0.1:5000/rede/dadosemarquivo/xlsx \
  -F "data=@grafo.json" \
  --output rede.xlsx
```

## ⚙️ Configuração Básica

Edite `rede.ini`:

```ini
[BASE]
# Seus bancos de dados
base_receita = bases/cnpj.db
base_rede = bases/rede.db

[ETC]
# Porta do servidor
porta_flask = 5000

# Rate limiting
limiter_padrao = 20/minute
```

## 🔧 Comandos Úteis

### Desenvolvimento

```bash
# Compilar
make build

# Executar testes
make test

# Verificar código
make vet

# Formatar código
make fmt

# Limpar arquivos gerados
make clean
```

### Produção

```bash
# Build otimizado
make build-prod

# Build para múltiplas plataformas
make build-all

# Executar em background
nohup ./rede-cnpj > app.log 2>&1 &
```

## 🐛 Troubleshooting

### Erro: "cannot find package"

```bash
go mod download
go mod tidy
make build
```

### Erro: "database not found"

```bash
# Verificar se os arquivos existem
ls -la bases/

# Copiar novamente se necessário
cp ../rede/bases/*.db bases/
```

### Erro: "port already in use"

```bash
# Usar outra porta
./rede-cnpj -porta_flask=8080
```

### Erro: "permission denied"

```bash
# Dar permissão de execução
chmod +x rede-cnpj
```

## 📊 Comparar com Python

### Python
```bash
cd ../rede
python rede.py
```

### Go
```bash
cd ../RedeGO
./rede-cnpj
```

Compare:
- Tempo de startup
- Uso de memória
- Velocidade de resposta

## 🎓 Próximos Passos

1. **Ler documentação completa**
   - [README.md](README.md) - Visão geral
   - [INSTALL.md](INSTALL.md) - Instalação detalhada
   - [ARCHITECTURE.md](ARCHITECTURE.md) - Arquitetura

2. **Explorar funcionalidades**
   - Buscar empresas
   - Visualizar relacionamentos
   - Exportar dados

3. **Personalizar configuração**
   - Editar `rede.ini`
   - Ajustar limites
   - Configurar rate limiting

4. **Testar com dados reais**
   - Copiar bancos completos
   - Fazer buscas complexas
   - Medir performance

5. **Contribuir**
   - Reportar bugs
   - Sugerir melhorias
   - Enviar pull requests

## 📚 Recursos

- **Documentação**: Ver arquivos `.md` no projeto
- **Exemplos**: Ver `pkg/cpfcnpj/validator_test.go`
- **Issues**: https://github.com/peder1981/rede-cnpj/issues

## 💡 Dicas

1. **Use make** para comandos comuns
2. **Leia os logs** se algo der errado
3. **Compare resultados** com versão Python
4. **Teste incrementalmente** antes de produção
5. **Faça backup** dos bancos de dados

## ✨ Funcionalidades Principais

- ✅ Busca por CNPJ/CPF
- ✅ Busca por nome
- ✅ Rede de relacionamentos
- ✅ Exportação Excel/JSON
- ✅ API RESTful
- ✅ Rate limiting
- ✅ Alta performance

## 🎉 Pronto!

Você agora tem o RedeCNPJ em Go funcionando!

Para mais informações, consulte:
- [SUMMARY.md](SUMMARY.md) - Resumo completo
- [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) - Migração Python→Go

---

**Dúvidas?** Abra uma issue no GitHub!
