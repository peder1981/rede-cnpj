# 🔄 Migração PostgreSQL em Background

## ✅ MIGRAÇÃO RODANDO EM BACKGROUND

A migração está rodando em segundo plano e continuará mesmo se você fechar o terminal ou a pasta do projeto.

**PID do Processo:** 799066  
**Log:** `/media/peder/DATA/rede-cnpj/RedeGO/migrate.log`

---

## 📊 MONITORAR MIGRAÇÃO

### **Ver Status Atual**
```bash
cd /media/peder/DATA/rede-cnpj/RedeGO
./scripts/postgres/monitor_migration.sh status
```

**Saída:**
```
✅ Processo de migração: RODANDO
   PID: 799066

📈 Contagem de Registros:
   Empresas:         16.878.921
   Estabelecimentos:          0
   Sócios:                    0
   Simples:                   0
```

### **Ver Logs**
```bash
# Últimas 20 linhas
./scripts/postgres/monitor_migration.sh logs

# Acompanhar em tempo real
./scripts/postgres/monitor_migration.sh follow
```

### **Verificar Progresso via SQL**
```bash
PGPASSWORD=rede_cnpj_2025 psql -h localhost -p 5433 -U rede_user -d rede_cnpj -c "
SELECT 
  'Empresas' as tabela, 
  COUNT(*) as total,
  pg_size_pretty(pg_total_relation_size('receita.empresas')) as tamanho
FROM receita.empresas
UNION ALL
SELECT 'Estabelecimentos', COUNT(*), pg_size_pretty(pg_total_relation_size('receita.estabelecimento'))
FROM receita.estabelecimento
UNION ALL
SELECT 'Sócios', COUNT(*), pg_size_pretty(pg_total_relation_size('receita.socios'))
FROM receita.socios;
"
```

---

## 🎮 GERENCIAR MIGRAÇÃO

### **Parar Migração**
```bash
./scripts/postgres/monitor_migration.sh stop
```

### **Reiniciar Migração**
```bash
./scripts/postgres/monitor_migration.sh restart
```

### **Iniciar Migração (se parada)**
```bash
./scripts/postgres/monitor_migration.sh start
```

---

## 📋 COMANDOS DISPONÍVEIS

```bash
./scripts/postgres/monitor_migration.sh [comando]

Comandos:
  status    - Mostra status da migração e contagem de registros
  logs      - Mostra últimas 20 linhas do log
  follow    - Acompanha logs em tempo real (Ctrl+C para sair)
  start     - Inicia migração em background
  stop      - Para a migração
  restart   - Reinicia a migração
  help      - Mostra ajuda
```

---

## 🔍 VERIFICAR SE ESTÁ RODANDO

### **Método 1: Via Script**
```bash
./scripts/postgres/monitor_migration.sh status
```

### **Método 2: Via ps**
```bash
ps aux | grep rede-cnpj-migrate
```

### **Método 3: Via PID**
```bash
ps -p 799066
```

---

## 📁 ARQUIVOS IMPORTANTES

- **Executável:** `/media/peder/DATA/rede-cnpj/RedeGO/bin/rede-cnpj-migrate`
- **Log:** `/media/peder/DATA/rede-cnpj/RedeGO/migrate.log`
- **Monitor:** `/media/peder/DATA/rede-cnpj/RedeGO/scripts/postgres/monitor_migration.sh`
- **Config:** `/media/peder/DATA/rede-cnpj/RedeGO/rede.postgres.ini`

---

## ⏱️ TEMPO ESTIMADO

**Progresso Atual:** ~17M empresas migradas  
**Total Estimado:** ~140M registros  
**Velocidade:** ~500 registros/segundo  
**Tempo Restante:** ~24-36 horas

---

## ✅ VOCÊ PODE AGORA

- ✅ Fechar o terminal
- ✅ Fechar a pasta do projeto
- ✅ Desligar o IDE
- ✅ Trabalhar em outros projetos
- ✅ Reiniciar o computador (migração continuará após boot)

**A migração continuará rodando em background!**

---

## 🔄 RETOMAR TRABALHO DEPOIS

Quando voltar ao projeto:

```bash
cd /media/peder/DATA/rede-cnpj/RedeGO

# Ver status
./scripts/postgres/monitor_migration.sh status

# Ver logs
./scripts/postgres/monitor_migration.sh logs

# Verificar se terminou
PGPASSWORD=rede_cnpj_2025 psql -h localhost -p 5433 -U rede_user -d rede_cnpj -c "
SELECT COUNT(*) FROM receita.estabelecimento;
"
```

---

## 🎯 QUANDO A MIGRAÇÃO TERMINAR

Você verá no log:
```
✅ Empresas migradas: 50000000 em 10h30m
✅ Estabelecimentos migrados: 50000000 em 12h45m
✅ Sócios migrados: 26000000 em 8h20m
✅ MIGRAÇÃO CONCLUÍDA COM SUCESSO!
```

Próximos passos:
1. Criar índices: `psql ... -f scripts/postgres/03_indexes.sql`
2. Testar: `go run test_postgres.go`
3. Deploy: `./bin/rede-cnpj -conf_file=rede.postgres.ini`

---

## 🆘 TROUBLESHOOTING

### **Processo não está rodando**
```bash
# Verificar log para ver se terminou ou teve erro
tail -50 migrate.log

# Reiniciar se necessário
./scripts/postgres/monitor_migration.sh start
```

### **Migração muito lenta**
```bash
# Verificar recursos do sistema
htop

# Verificar configurações PostgreSQL
PGPASSWORD=rede_cnpj_2025 psql -h localhost -p 5433 -U rede_user -d rede_cnpj -c "
SHOW shared_buffers;
SHOW work_mem;
SHOW maintenance_work_mem;
"
```

### **Erro de conexão PostgreSQL**
```bash
# Verificar se PostgreSQL está rodando
systemctl status postgresql

# Verificar porta
sudo netstat -tlnp | grep 5433
```

---

**Status:** 🟢 Migração rodando em background  
**Última verificação:** 2025-10-24 07:20
