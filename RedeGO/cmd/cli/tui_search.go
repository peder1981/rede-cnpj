package main

import (
	"fmt"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) viewSearch() string {
	s := "\n"
	s += "╔══════════════════════════════════════════════════════════════════════╗\n"
	s += "║         🔍 Buscar CPF ou CNPJ                                        ║\n"
	s += "╚══════════════════════════════════════════════════════════════════════╝\n\n"

	s += "Digite o CPF (11 dígitos) ou CNPJ (14 dígitos) sem pontuação:\n\n"
	s += fmt.Sprintf("┌──────────────────────────────────────────────────────────────────────┐\n")
	s += fmt.Sprintf("│ > %-67s│\n", m.searchInput)
	s += fmt.Sprintf("└──────────────────────────────────────────────────────────────────────┘\n\n")

	s += "Exemplos:\n"
	s += "  CPF:  12345678900\n"
	s += "  CNPJ: 01234567000100\n\n"

	s += "┌──────────────────────────────────────────────────────────────────────┐\n"
	s += "│ [ENTER] Buscar | [Q] Cancelar | [ESC] Voltar                        │\n"
	s += "│                                                                      │\n"
	s += "│ Após buscar:                                                         │\n"
	s += "│ • CPF:  Ver todas empresas + Investigação forense                   │\n"
	s += "│ • CNPJ: Ver dados completos + Sócios                                │\n"
	s += "└──────────────────────────────────────────────────────────────────────┘\n"

	if m.message != "" {
		s += "\n"
		s += fmt.Sprintf("💬 %s\n", m.message)
	}

	return s
}

func (m model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Valida e busca
		cleaned := cleanInput(m.searchInput)
		
		if len(cleaned) == 11 {
			// CPF
			m.viewData = cleaned
			m.mode = modeViewCPF
			m.message = "Exibindo dados do CPF"
		} else if len(cleaned) == 14 {
			// CNPJ
			m.viewData = cleaned
			m.mode = modeViewCNPJ
			m.message = "Exibindo dados do CNPJ"
		} else {
			m.message = "❌ Formato inválido! Use 11 dígitos (CPF) ou 14 (CNPJ)"
		}
		
	case "backspace":
		if len(m.searchInput) > 0 {
			m.searchInput = m.searchInput[:len(m.searchInput)-1]
		}
		
	case "q":
		m.mode = modeTree
		m.searchInput = ""
		m.message = "Busca cancelada"
		
	default:
		// Aceita apenas números
		if len(msg.String()) == 1 && msg.String() >= "0" && msg.String() <= "9" {
			if len(m.searchInput) < 14 {
				m.searchInput += msg.String()
			}
		}
	}
	
	return m, nil
}

func cleanInput(s string) string {
	// Remove tudo que não é número
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(s, "")
}

// updateViewCPF atualiza visualização de CPF
func (m model) updateViewCPF(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "backspace":
		m.mode = modeTree
		m.message = "Voltou ao modo árvore"
	case "i":
		// Investigação forense
		m.mode = modeForensics
		m.message = "Análise forense iniciada"
	case "t":
		// Timeline
		m.mode = modeTimeline
		m.message = "Carregando timeline de atividades..."
	}
	return m, nil
}

// updateViewCNPJ atualiza visualização de CNPJ
func (m model) updateViewCNPJ(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "backspace":
		m.mode = modeTree
		m.message = "Voltou ao modo árvore"
	case "s":
		// Ver sócios
		m.mode = modeViewSocios
		m.message = "Carregando lista de sócios..."
	case "c":
		// Cadeia de controle
		m.mode = modeCadeiaControle
		m.message = "Carregando cadeia de controle..."
	}
	return m, nil
}

// updateForensics atualiza visualização forense
func (m model) updateForensics(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "backspace":
		m.mode = modeViewCPF
		m.message = "Voltou para dados do CPF"
	case "d":
		// Ver detalhes de empresa - precisa selecionar empresa primeiro
		m.mode = modeEmpresaDetalhes
		m.selectedEmpresaCursor = 0
		m.message = "Selecione uma empresa para ver detalhes"
	case "t":
		// Timeline
		m.mode = modeTimeline
		m.message = "Carregando timeline de atividades..."
	}
	return m, nil
}

// updateViewSocios atualiza visualização de sócios
func (m model) updateViewSocios(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "backspace":
		m.mode = modeViewCNPJ
		m.message = "Voltou para dados do CNPJ"
	}
	return m, nil
}

// updateCadeiaControle atualiza visualização de cadeia de controle
func (m model) updateCadeiaControle(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "backspace":
		m.mode = modeViewCNPJ
		m.message = "Voltou para dados do CNPJ"
	}
	return m, nil
}

// updateTimeline atualiza visualização de timeline
func (m model) updateTimeline(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "backspace":
		// Volta para o modo anterior (CPF ou Forensics)
		if m.searchType == "cpf" {
			m.mode = modeViewCPF
		} else {
			m.mode = modeForensics
		}
		m.message = "Voltou para visualização anterior"
	}
	return m, nil
}

// updateEmpresaDetalhes atualiza visualização de detalhes de empresa
func (m model) updateEmpresaDetalhes(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "backspace":
		m.mode = modeForensics
		m.message = "Voltou para análise forense"
	}
	return m, nil
}
