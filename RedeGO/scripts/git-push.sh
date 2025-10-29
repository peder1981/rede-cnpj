#!/bin/bash

# Script para fazer push com Bagus Browser configurado
# Uso: ./scripts/git-push.sh [mensagem do commit]

# Configurar Bagus Browser
export BROWSER=/usr/local/bin/bagus-browser

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Git Push com Bagus Browser${NC}"
echo ""

# Verificar se há mudanças
if [[ -z $(git status -s) ]]; then
    echo -e "${YELLOW}⚠️  Nenhuma mudança para commitar${NC}"
    echo ""
    echo "Verificando commits pendentes..."
    COMMITS_AHEAD=$(git rev-list --count @{u}..HEAD 2>/dev/null || echo "0")
    
    if [ "$COMMITS_AHEAD" -gt 0 ]; then
        echo -e "${YELLOW}📤 Você tem $COMMITS_AHEAD commit(s) para enviar${NC}"
        echo ""
        read -p "Deseja fazer push agora? (s/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Ss]$ ]]; then
            echo -e "${GREEN}📤 Fazendo push...${NC}"
            git push origin master
            exit $?
        else
            echo -e "${YELLOW}Push cancelado${NC}"
            exit 0
        fi
    else
        echo -e "${GREEN}✓ Tudo sincronizado!${NC}"
        exit 0
    fi
fi

# Adicionar todas as mudanças
echo -e "${GREEN}📋 Adicionando arquivos...${NC}"
git add -A

# Mostrar status
echo ""
git status --short

# Pedir mensagem de commit se não foi fornecida
if [ -z "$1" ]; then
    echo ""
    echo -e "${YELLOW}💬 Digite a mensagem do commit:${NC}"
    read -r COMMIT_MSG
else
    COMMIT_MSG="$*"
fi

# Fazer commit
echo ""
echo -e "${GREEN}💾 Fazendo commit...${NC}"
git commit -m "$COMMIT_MSG"

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Erro ao fazer commit${NC}"
    exit 1
fi

# Fazer push
echo ""
echo -e "${GREEN}📤 Fazendo push...${NC}"
git push origin master

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Push realizado com sucesso!${NC}"
else
    echo ""
    echo -e "${RED}❌ Erro ao fazer push${NC}"
    exit 1
fi
