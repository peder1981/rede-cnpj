package main

import (
	"fmt"
	"strings"
)

func (m model) viewCrossData() string {
	s := "\n"
	s += "╔══════════════════════════════════════════════════════════════════════╗\n"
	s += "║         🔓 RedeCNPJ - Cruzamento de Dados (SEM CENSURA)             ║\n"
	s += "╚══════════════════════════════════════════════════════════════════════╝\n\n"

	s += "Selecione o tipo de cruzamento:\n\n"

	options := []string{
		"1. 📋 CPF → Empresas - Todas as empresas de uma pessoa",
		"2. 👥 CNPJ → Sócios - Todos os sócios (CPF completo)",
		"3. 🔗 Sócios em Comum - Entre duas empresas",
		"4. 🕸️  Rede 2º Grau - Empresas dos sócios",
		"5. 🏠 Mesmo Endereço - Empresas no mesmo local",
		"6. 📞 Mesmo Contato - Email/telefone compartilhado",
		"7. 👶 Representantes Legais - Menores + representantes",
		"8. 🌍 Empresas Estrangeiras - Sede no exterior",
		"9. 🌎 Sócios Estrangeiros - Pessoas estrangeiras",
		"10. 📅 Timeline - Histórico completo de atividades",
		"11. ⚠️  Empresas Baixadas - Padrões de comportamento",
		"12. 📊 Dados Completos - TUDO sem censura",
	}

	for i, opt := range options {
		cursor := "  "
		if m.crossMenu == i {
			cursor = "→ "
		}
		s += fmt.Sprintf("%s%s\n", cursor, opt)
	}

	s += "\n"
	s += "┌──────────────────────────────────────────────────────────────────────┐\n"
	s += "│ ⚠️  ATENÇÃO: Dados exibidos SEM CENSURA                              │\n"
	s += "│ • CPF completo (sem máscara)                                         │\n"
	s += "│ • Endereços completos                                                │\n"
	s += "│ • Telefones e emails                                                 │\n"
	s += "│ • Dados de empresários individuais                                   │\n"
	s += "│                                                                      │\n"
	s += "│ Use com responsabilidade conforme LGPD                               │\n"
	s += "└──────────────────────────────────────────────────────────────────────┘\n\n"

	s += "┌──────────────────────────────────────────────────────────────────────┐\n"
	s += "│ ↑↓ Navegar | [ENTER] Selecionar | [Q] Voltar                        │\n"
	s += "│                                                                      │\n"
	s += "│ NOTA: Funcionalidade de cruzamento disponível via API REST          │\n"
	s += "│ Consulte CROSSDATA_API.md para exemplos de uso                      │\n"
	s += "└──────────────────────────────────────────────────────────────────────┘\n"

	if m.message != "" {
		s += "\n"
		s += "┌──────────────────────────────────────────────────────────────────────┐\n"
		s += fmt.Sprintf("│ %s%s │\n", m.message, strings.Repeat(" ", 68-len(m.message)))
		s += "└──────────────────────────────────────────────────────────────────────┘\n"
	}

	return s
}
