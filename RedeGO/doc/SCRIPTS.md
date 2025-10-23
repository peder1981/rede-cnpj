# Scripts de Gerenciamento - RedeCNPJ Go

Documentação dos scripts shell para gerenciar a aplicação RedeCNPJ.

## 📜 Scripts Disponíveis

### 1. `start.sh` - Inicialização Completa

Script principal que realiza todas as tarefas necessárias e inicia a aplicação.

#### Funcionalidades

✅ **Verificações Automáticas**
- Verifica se Go está instalado
- Valida versão do Go (1.21+)
- Verifica estrutura de diretórios
- Checa existência de bancos de dados

✅ **Preparação do Ambiente**
- Cria diretórios necessários (bases/, arquivos/, logs/)
- Instala dependências Go
- Compila a aplicação
- Valida configuração

✅ **Gerenciamento de Processos**
- Inicia o servidor
- Captura CTRL+C para encerramento gracioso
- Encerra todos os processos filhos ao sair
- Aguarda até 5s para shutdown gracioso
- Force kill se necessário

✅ **Logs**
- Output colorido no terminal
- Log completo em `logs/server.log`
- Mensagens informativas de progresso

#### Uso Básico

```bash
# Iniciar com configurações padrão
./start.sh

# Iniciar na porta 8080
./start.sh -p 8080

# Usar arquivo de configuração customizado
./start.sh -c custom.ini

# Modo desenvolvimento (recompila ao detectar mudanças)
./start.sh -d
```

#### Opções

| Opção | Descrição | Padrão |
|-------|-----------|--------|
| `-p, --port PORT` | Porta do servidor | 5000 |
| `-c, --config FILE` | Arquivo de configuração | rede.ini |
| `-d, --dev` | Modo desenvolvimento | false |
| `-h, --help` | Exibe ajuda | - |

#### Exemplos

```bash
# Desenvolvimento com porta customizada
./start.sh -d -p 8080

# Produção com configuração específica
./start.sh -c production.ini -p 80

# Apenas ajuda
./start.sh --help
```

#### Encerramento

Para encerrar a aplicação, pressione **CTRL+C**. O script irá:

1. Capturar o sinal de interrupção
2. Encerrar o servidor graciosamente (SIGTERM)
3. Aguardar até 5 segundos
4. Forçar encerramento se necessário (SIGKILL)
5. Limpar processos filhos
6. Exibir mensagem de confirmação

### 2. `stop.sh` - Encerramento Manual

Script para encerrar todos os processos do RedeCNPJ em execução.

#### Uso

```bash
./stop.sh
```

#### Funcionalidades

- Procura todos os processos `rede-cnpj`
- Encerra graciosamente (SIGTERM)
- Aguarda 2 segundos
- Force kill processos restantes (SIGKILL)
- Confirma encerramento

#### Quando Usar

- Quando `start.sh` foi encerrado abruptamente
- Para limpar processos órfãos
- Em scripts de automação
- Antes de fazer manutenção

### 3. `restart.sh` - Reinicialização

Script para reiniciar a aplicação (encerra e inicia novamente).

#### Uso

```bash
# Reiniciar com configurações padrão
./restart.sh

# Reiniciar com opções customizadas
./restart.sh -p 8080 -d
```

#### Funcionalidades

- Executa `stop.sh` para encerrar processos
- Aguarda 2 segundos
- Executa `start.sh` com as opções fornecidas

#### Quando Usar

- Após mudanças na configuração
- Após atualização do código
- Para aplicar novas variáveis de ambiente
- Em rotinas de manutenção

## 🎯 Fluxo de Execução do start.sh

```
┌─────────────────────────────────────┐
│  1. Verificar Pré-requisitos        │
│     - Go instalado?                 │
│     - Versão correta?               │
│     - Diretório correto?            │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  2. Criar Diretórios                │
│     - bases/                        │
│     - arquivos/                     │
│     - static/                       │
│     - templates/                    │
│     - logs/                         │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  3. Verificar Bancos de Dados       │
│     - Existem em bases/?            │
│     - Copiar de ../rede/bases/?     │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  4. Instalar Dependências           │
│     - go mod download               │
│     - go mod tidy                   │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  5. Compilar Aplicação              │
│     - go build -o rede-cnpj         │
│     - chmod +x rede-cnpj            │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  6. Verificar Configuração          │
│     - rede.ini existe?              │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  7. Iniciar Servidor                │
│     - ./rede-cnpj [opções]          │
│     - Capturar PID                  │
│     - Registrar trap SIGINT/SIGTERM │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  8. Aguardar CTRL+C                 │
│     - Servidor rodando...           │
└──────────────┬──────────────────────┘
               │
         [CTRL+C pressionado]
               │
┌──────────────▼──────────────────────┐
│  9. Cleanup (função trap)           │
│     - kill -TERM $APP_PID           │
│     - Aguardar até 5s               │
│     - kill -9 se necessário         │
│     - Encerrar processos filhos     │
│     - Limpar recursos               │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  10. Sair                           │
│      - Mensagem de confirmação      │
│      - exit 0                       │
└─────────────────────────────────────┘
```

## 🔧 Modo Desenvolvimento

O modo desenvolvimento (`-d`) oferece recursos adicionais:

### Funcionalidades

- **Auto-reload**: Recompila e reinicia ao detectar mudanças
- **Watch de arquivos**: Monitora todos os arquivos `.go`
- **Feedback imediato**: Mostra quando recompilação ocorre

### Requisitos

Para usar o modo desenvolvimento com watch automático:

```bash
# Ubuntu/Debian
sudo apt-get install inotify-tools

# Fedora/RHEL
sudo dnf install inotify-tools

# Arch Linux
sudo pacman -S inotify-tools
```

### Uso

```bash
./start.sh -d
```

Agora, qualquer mudança em arquivos `.go` irá:
1. Detectar a mudança
2. Recompilar a aplicação
3. Reiniciar o servidor automaticamente

## 📝 Logs

### Localização

- **Terminal**: Output colorido em tempo real
- **Arquivo**: `logs/server.log`

### Formato

```
[INFO] Mensagens informativas (verde)
[WARN] Avisos (amarelo)
[ERRO] Erros (vermelho)
[STEP] Passos do processo (azul)
```

### Visualizar Logs

```bash
# Tempo real
tail -f logs/server.log

# Últimas 100 linhas
tail -n 100 logs/server.log

# Buscar erros
grep ERROR logs/server.log

# Limpar logs antigos
> logs/server.log
```

## 🛡️ Tratamento de Sinais

O script captura e trata os seguintes sinais:

| Sinal | Descrição | Ação |
|-------|-----------|------|
| SIGINT | CTRL+C | Cleanup e saída gracioso |
| SIGTERM | Kill normal | Cleanup e saída gracioso |
| EXIT | Qualquer saída | Cleanup automático |

### Sequência de Encerramento

1. **Sinal recebido** → Trap ativado
2. **SIGTERM enviado** → Servidor recebe sinal de encerramento
3. **Aguarda 5s** → Tempo para cleanup gracioso
4. **Verifica processo** → Ainda está rodando?
5. **SIGKILL se necessário** → Force kill
6. **Limpa filhos** → Encerra processos filhos
7. **Mensagem final** → Confirma encerramento

## 🔍 Troubleshooting

### Problema: "Go não está instalado"

```bash
# Instalar Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### Problema: "Porta já em uso"

```bash
# Ver o que está usando a porta
sudo lsof -i :5000

# Usar outra porta
./start.sh -p 8080
```

### Problema: "Falha na compilação"

```bash
# Limpar e recompilar
make clean
./start.sh
```

### Problema: "Bancos de dados não encontrados"

```bash
# Copiar da versão Python
cp ../rede/bases/*.db bases/

# Ou criar links simbólicos
ln -s ../rede/bases/cnpj_teste.db bases/
ln -s ../rede/bases/rede_teste.db bases/
```

### Problema: "Processo não encerra com CTRL+C"

```bash
# Usar script de stop
./stop.sh

# Ou manualmente
pkill -9 rede-cnpj
```

## 🚀 Uso em Produção

### Systemd Service

Para rodar como serviço do sistema:

```bash
# Criar arquivo de serviço
sudo nano /etc/systemd/system/rede-cnpj.service
```

Conteúdo:

```ini
[Unit]
Description=RedeCNPJ Go Server
After=network.target

[Service]
Type=simple
User=seu-usuario
WorkingDirectory=/caminho/para/RedeGO
ExecStart=/caminho/para/RedeGO/start.sh -p 5000
ExecStop=/caminho/para/RedeGO/stop.sh
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Ativar:

```bash
sudo systemctl daemon-reload
sudo systemctl enable rede-cnpj
sudo systemctl start rede-cnpj
sudo systemctl status rede-cnpj
```

### Cron Job (Reinicialização Automática)

```bash
# Editar crontab
crontab -e

# Reiniciar diariamente às 3h
0 3 * * * /caminho/para/RedeGO/restart.sh >> /var/log/rede-cnpj-cron.log 2>&1
```

## 📊 Exemplos de Uso

### Desenvolvimento Local

```bash
# Terminal 1: Servidor em modo dev
./start.sh -d

# Terminal 2: Fazer mudanças no código
vim internal/handlers/handlers.go
# Salvar → Auto-reload acontece

# Terminal 3: Testar
curl http://localhost:5000/rede/api/status
```

### Teste de Carga

```bash
# Terminal 1: Iniciar servidor
./start.sh

# Terminal 2: Apache Bench
ab -n 1000 -c 10 http://localhost:5000/rede/api/status

# Terminal 1: CTRL+C para encerrar
```

### Deploy Rápido

```bash
# Atualizar código
git pull

# Reiniciar
./restart.sh -c production.ini

# Verificar
curl http://localhost:5000/rede/api/status
```

## ✅ Checklist de Uso

Antes de executar `start.sh`:

- [ ] Go 1.21+ instalado
- [ ] Bancos de dados em `bases/`
- [ ] Arquivo `rede.ini` configurado
- [ ] Porta disponível (padrão: 5000)
- [ ] Permissões de execução (`chmod +x start.sh`)

Após executar:

- [ ] Servidor iniciou sem erros
- [ ] Acessível em http://localhost:5000/rede/
- [ ] Logs sendo gerados em `logs/server.log`
- [ ] CTRL+C encerra graciosamente

## 🎓 Dicas

1. **Use modo dev** durante desenvolvimento: `./start.sh -d`
2. **Monitore logs** em tempo real: `tail -f logs/server.log`
3. **Teste encerramento** antes de produção
4. **Configure systemd** para produção
5. **Faça backup** dos bancos de dados regularmente

## 📞 Suporte

Se encontrar problemas:

1. Verifique logs em `logs/server.log`
2. Execute `./start.sh -h` para ver opções
3. Teste com `./stop.sh` e `./start.sh`
4. Reporte issues no GitHub

---

**Scripts criados para facilitar o gerenciamento do RedeCNPJ Go!**
