package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/analytics"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/export"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/models"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/services"
)

type viewMode int

const (
	modeTree viewMode = iota
	modeAnalytics
	modeHelp
	modeExport
)

type nodeItem struct {
	node     models.Node
	level    int
	expanded bool
	parent   string
}

type model struct {
	redeService *services.RedeService
	rootCNPJ    string
	items       []nodeItem
	cursor      int
	width       int
	height      int
	err         error
	mode        viewMode
	message     string
	stats       *analytics.GraphStats
	currentGraph *models.Graph
	exportMenu   int // 0=excel, 1=csv_nodes, 2=csv_edges, 3=csv_stats
}

func initialModel(redeService *services.RedeService, cnpj string) model {
	return model{
		redeService: redeService,
		rootCNPJ:    cnpj,
		items:       []nodeItem{},
		cursor:      0,
		mode:        modeTree,
		exportMenu:  0,
	}
}

func (m model) Init() tea.Cmd {
	return m.loadRoot()
}

func (m model) loadRoot() tea.Cmd {
	return func() tea.Msg {
		graph, err := m.redeService.CamadasRede(1, []string{m.rootCNPJ}, "", "")
		if err != nil {
			return errMsg{err}
		}
		return graphMsg{graph}
	}
}

func (m model) expandNode(nodeID string) tea.Cmd {
	return func() tea.Msg {
		graph, err := m.redeService.CamadasRede(1, []string{nodeID}, "", "")
		if err != nil {
			return errMsg{err}
		}
		return expandMsg{nodeID, graph}
	}
}

type graphMsg struct {
	graph *models.Graph
}

type expandMsg struct {
	parentID string
	graph    *models.Graph
}

type errMsg struct {
	err error
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Comandos globais
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "F1", "?":
			m.mode = modeHelp
			return m, nil
		}

		// Comandos específicos por modo
		switch m.mode {
		case modeTree:
			return m.updateTree(msg)
		case modeAnalytics:
			return m.updateAnalytics(msg)
		case modeHelp:
			return m.updateHelp(msg)
		case modeExport:
			return m.updateExport(msg)
		}

	case graphMsg:
		m.currentGraph = msg.graph
		m.items = []nodeItem{}
		for _, node := range msg.graph.Nodes {
			m.items = append(m.items, nodeItem{
				node:     node,
				level:    0,
				expanded: false,
				parent:   "",
			})
		}
		m.message = fmt.Sprintf("✓ Carregado: %d nós, %d ligações", len(msg.graph.Nodes), len(msg.graph.Edges))

	case expandMsg:
		parentIdx := -1
		for i, item := range m.items {
			if item.node.ID == msg.parentID {
				parentIdx = i
				break
			}
		}

		if parentIdx >= 0 {
			parentLevel := m.items[parentIdx].level
			m.items = m.collapseChildren(parentIdx)
			
			newItems := []nodeItem{}
			for _, node := range msg.graph.Nodes {
				if node.ID != msg.parentID {
					newItems = append(newItems, nodeItem{
						node:     node,
						level:    parentLevel + 1,
						expanded: false,
						parent:   msg.parentID,
					})
				}
			}

			m.items = append(
				m.items[:parentIdx+1],
				append(newItems, m.items[parentIdx+1:]...)...,
			)
			
			// Atualiza grafo atual
			if m.currentGraph != nil {
				m.currentGraph.Nodes = append(m.currentGraph.Nodes, msg.graph.Nodes...)
				m.currentGraph.Edges = append(m.currentGraph.Edges, msg.graph.Edges...)
			}
			
			m.message = fmt.Sprintf("✓ Expandido: +%d nós", len(msg.graph.Nodes)-1)
		}

	case errMsg:
		m.err = msg.err
		m.message = fmt.Sprintf("✗ Erro: %v", msg.err)
	}

	return m, nil
}

func (m model) updateTree(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}
	case "enter", " ", "right", "l":
		if len(m.items) > 0 && m.cursor < len(m.items) {
			item := &m.items[m.cursor]
			if !item.expanded {
				item.expanded = true
				m.message = "⏳ Expandindo..."
				return m, m.expandNode(item.node.ID)
			}
		}
	case "left", "h":
		if len(m.items) > 0 && m.cursor < len(m.items) {
			item := &m.items[m.cursor]
			if item.expanded {
				item.expanded = false
				m.items = m.collapseChildren(m.cursor)
				m.message = "✓ Colapsado"
			}
		}
	case "a":
		if m.currentGraph != nil && len(m.currentGraph.Nodes) > 0 {
			m.mode = modeAnalytics
			analyzer := analytics.NewAnalyzer()
			m.stats = analyzer.AnalyzeGraph(m.currentGraph)
			m.message = "✓ Estatísticas calculadas"
		} else {
			m.message = "✗ Carregue um grafo primeiro"
		}
	case "e":
		if m.currentGraph != nil && len(m.currentGraph.Nodes) > 0 {
			m.mode = modeExport
			m.exportMenu = 0
			m.message = "Selecione o formato de exportação"
		} else {
			m.message = "✗ Carregue um grafo primeiro"
		}
	}
	return m, nil
}

func (m model) updateAnalytics(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "backspace":
		m.mode = modeTree
		m.message = "Voltou ao modo árvore"
	}
	return m, nil
}

func (m model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "backspace", "F1", "?":
		m.mode = modeTree
		m.message = "Voltou ao modo árvore"
	}
	return m, nil
}

func (m model) updateExport(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.exportMenu > 0 {
			m.exportMenu--
		}
	case "down", "j":
		if m.exportMenu < 3 {
			m.exportMenu++
		}
	case "enter", " ":
		// Executa exportação
		filename := m.doExport()
		if filename != "" {
			m.message = fmt.Sprintf("✓ Exportado: %s", filename)
		} else {
			m.message = "✗ Erro ao exportar"
		}
		m.mode = modeTree
	case "q", "backspace":
		m.mode = modeTree
		m.message = "Exportação cancelada"
	}
	return m, nil
}

func (m model) doExport() string {
	if m.currentGraph == nil {
		return ""
	}

	outputDir := "output"
	os.MkdirAll(outputDir, 0755)

	switch m.exportMenu {
	case 0: // Excel
		exporter := export.NewExcelExporter()
		defer exporter.Close()
		
		data, err := exporter.ExportGraph(m.currentGraph)
		if err != nil {
			return ""
		}
		
		filename := filepath.Join(outputDir, "rede-cnpj.xlsx")
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return ""
		}
		return filename

	case 1: // CSV Nós
		exporter := export.NewCSVExporter()
		data, err := exporter.ExportNodes(m.currentGraph.Nodes)
		if err != nil {
			return ""
		}
		
		filename := filepath.Join(outputDir, "nos.csv")
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return ""
		}
		return filename

	case 2: // CSV Arestas
		exporter := export.NewCSVExporter()
		data, err := exporter.ExportEdges(m.currentGraph.Edges)
		if err != nil {
			return ""
		}
		
		filename := filepath.Join(outputDir, "arestas.csv")
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return ""
		}
		return filename

	case 3: // CSV Estatísticas
		exporter := export.NewCSVExporter()
		data, err := exporter.ExportStats(m.currentGraph)
		if err != nil {
			return ""
		}
		
		filename := filepath.Join(outputDir, "estatisticas.csv")
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return ""
		}
		return filename
	}

	return ""
}

func (m model) collapseChildren(parentIdx int) []nodeItem {
	if parentIdx >= len(m.items) {
		return m.items
	}

	parentID := m.items[parentIdx].node.ID
	parentLevel := m.items[parentIdx].level

	newItems := []nodeItem{}
	skip := false

	for i, item := range m.items {
		if i == parentIdx {
			newItems = append(newItems, item)
			skip = true
			continue
		}

		if skip {
			if item.level > parentLevel && (item.parent == parentID || m.isDescendant(item, parentID)) {
				continue
			} else {
				skip = false
			}
		}

		newItems = append(newItems, item)
	}

	return newItems
}

func (m model) isDescendant(item nodeItem, ancestorID string) bool {
	current := item.parent
	for current != "" {
		if current == ancestorID {
			return true
		}
		for _, i := range m.items {
			if i.node.ID == current {
				current = i.parent
				break
			}
		}
	}
	return false
}

func (m model) View() string {
	if m.err != nil {
		return m.viewError()
	}

	switch m.mode {
	case modeTree:
		return m.viewTree()
	case modeAnalytics:
		return m.viewAnalytics()
	case modeHelp:
		return m.viewHelp()
	case modeExport:
		return m.viewExport()
	}

	return ""
}

func (m model) viewError() string {
	return fmt.Sprintf("\n❌ ERRO: %v\n\nPressione ESC para sair.\n", m.err)
}

func (m model) viewTree() string {
	s := "\n"
	s += "╔══════════════════════════════════════════════════════════════════════╗\n"
	s += "║         🔍 RedeCNPJ - Navegação em Árvore                           ║\n"
	s += "╚══════════════════════════════════════════════════════════════════════╝\n\n"

	if len(m.items) == 0 {
		s += "⏳ Carregando...\n"
	} else {
		// Mostra até 20 itens
		start := 0
		end := len(m.items)
		if end > 20 {
			if m.cursor > 10 {
				start = m.cursor - 10
			}
			end = start + 20
			if end > len(m.items) {
				end = len(m.items)
				start = end - 20
				if start < 0 {
					start = 0
				}
			}
		}

		for i := start; i < end; i++ {
			item := m.items[i]
			
			cursor := "  "
			if m.cursor == i {
				cursor = "→ "
			}

			indent := strings.Repeat("  ", item.level)
			
			icon := "🏢"
			if len(item.node.ID) >= 3 && (item.node.ID[:3] == "PF_" || item.node.ID[:3] == "PE_") {
				icon = "👤"
			}

			expandIcon := "▶"
			if item.expanded {
				expandIcon = "▼"
			}

			label := item.node.Label
			if len(label) > 50 {
				label = label[:47] + "..."
			}

			s += fmt.Sprintf("%s%s%s %s %s\n", cursor, indent, expandIcon, icon, label)
		}
		
		if len(m.items) > 20 {
			s += fmt.Sprintf("\n[Mostrando %d-%d de %d itens]\n", start+1, end, len(m.items))
		}
	}

	s += "\n"
	
	// Barra de status
	if m.message != "" {
		s += "┌──────────────────────────────────────────────────────────────────────┐\n"
		s += fmt.Sprintf("│ %s%s │\n", m.message, strings.Repeat(" ", 68-len(m.message)))
		s += "└──────────────────────────────────────────────────────────────────────┘\n"
	}
	
	// Menu de comandos
	s += "┌──────────────────────────────────────────────────────────────────────┐\n"
	s += "│ NAVEGAÇÃO: ↑↓ mover | → expandir | ← colapsar                       │\n"
	s += "│ AÇÕES: [A]nalytics | [E]xportar | [F1/?] Ajuda | [ESC] Sair        │\n"
	s += "└──────────────────────────────────────────────────────────────────────┘\n"

	return s
}

func (m model) viewAnalytics() string {
	s := "\n"
	s += "╔══════════════════════════════════════════════════════════════════════╗\n"
	s += "║         📊 RedeCNPJ - Estatísticas do Grafo                         ║\n"
	s += "╚══════════════════════════════════════════════════════════════════════╝\n\n"

	if m.stats == nil {
		s += "Nenhuma estatística disponível.\n"
	} else {
		s += "┌─ RESUMO GERAL ─────────────────────────────────────────────────────┐\n"
		s += fmt.Sprintf("│ Total de Nós:           %6d                                     │\n", m.stats.TotalNodes)
		s += fmt.Sprintf("│ Total de Ligações:      %6d                                     │\n", m.stats.TotalEdges)
		s += fmt.Sprintf("│ Empresas (PJ):          %6d                                     │\n", m.stats.Empresas)
		s += fmt.Sprintf("│ Pessoas (PF):           %6d                                     │\n", m.stats.Pessoas)
		s += fmt.Sprintf("│ Pessoas Externas (PE):  %6d                                     │\n", m.stats.PessoasExternas)
		s += "└────────────────────────────────────────────────────────────────────┘\n\n"

		s += "┌─ MÉTRICAS DE REDE ─────────────────────────────────────────────────┐\n"
		s += fmt.Sprintf("│ Densidade:              %.4f                                   │\n", m.stats.Densidade)
		s += fmt.Sprintf("│ Grau Médio:             %.2f                                     │\n", m.stats.GrauMedio)
		s += fmt.Sprintf("│ Componentes Conexos:    %6d                                     │\n", m.stats.ComponentesConexos)
		s += "└────────────────────────────────────────────────────────────────────┘\n\n"

		s += "┌─ NÓS MAIS CONECTADOS ──────────────────────────────────────────────┐\n"
		for i, node := range m.stats.NosMaisConectados {
			if i >= 10 {
				break
			}
			label := node.Label
			if len(label) > 40 {
				label = label[:37] + "..."
			}
			s += fmt.Sprintf("│ %2d. %-40s [%3d conexões] │\n", i+1, label, node.Degree)
		}
		s += "└────────────────────────────────────────────────────────────────────┘\n\n"

		s += "┌─ TIPOS DE RELACIONAMENTO ──────────────────────────────────────────┐\n"
		count := 0
		for tipo, qtd := range m.stats.TiposRelacao {
			if count >= 10 {
				break
			}
			tipoLabel := tipo
			if len(tipoLabel) > 40 {
				tipoLabel = tipoLabel[:37] + "..."
			}
			s += fmt.Sprintf("│ %-40s: %6d                      │\n", tipoLabel, qtd)
			count++
		}
		s += "└────────────────────────────────────────────────────────────────────┘\n"
	}

	s += "\n"
	s += "┌──────────────────────────────────────────────────────────────────────┐\n"
	s += "│ Pressione [Q] ou [BACKSPACE] para voltar                            │\n"
	s += "└──────────────────────────────────────────────────────────────────────┘\n"

	return s
}

func (m model) viewExport() string {
	s := "\n"
	s += "╔══════════════════════════════════════════════════════════════════════╗\n"
	s += "║         💾 RedeCNPJ - Exportar Dados                                ║\n"
	s += "╚══════════════════════════════════════════════════════════════════════╝\n\n"

	s += "Selecione o formato de exportação:\n\n"

	options := []string{
		"📊 Excel (XLSX) - Arquivo completo com 3 planilhas",
		"📄 CSV - Nós (lista de entidades)",
		"📄 CSV - Arestas (lista de relacionamentos)",
		"📄 CSV - Estatísticas (resumo do grafo)",
	}

	for i, opt := range options {
		cursor := "  "
		if m.exportMenu == i {
			cursor = "→ "
		}
		s += fmt.Sprintf("%s%s\n", cursor, opt)
	}

	s += "\n"
	s += "┌──────────────────────────────────────────────────────────────────────┐\n"
	s += "│ ↑↓ Navegar | [ENTER] Exportar | [Q] Cancelar                        │\n"
	s += "│ Arquivos serão salvos em: ./output/                                 │\n"
	s += "└──────────────────────────────────────────────────────────────────────┘\n"

	return s
}

func (m model) viewHelp() string {
	s := "\n"
	s += "╔══════════════════════════════════════════════════════════════════════╗\n"
	s += "║         ❓ RedeCNPJ - Ajuda                                          ║\n"
	s += "╚══════════════════════════════════════════════════════════════════════╝\n\n"

	s += "┌─ MODO ÁRVORE (Navegação) ──────────────────────────────────────────┐\n"
	s += "│                                                                      │\n"
	s += "│  ↑ / k          - Mover cursor para cima                            │\n"
	s += "│  ↓ / j          - Mover cursor para baixo                           │\n"
	s += "│  → / l / ENTER  - Expandir nó selecionado (busca relacionamentos)   │\n"
	s += "│  ← / h          - Colapsar nó selecionado                           │\n"
	s += "│  a              - Ver Analytics (estatísticas do grafo)             │\n"
	s += "│  e              - Exportar dados (Excel ou CSV)                     │\n"
	s += "│  F1 / ?         - Mostrar esta ajuda                                │\n"
	s += "│  ESC            - Sair do programa                                  │\n"
	s += "│                                                                      │\n"
	s += "└──────────────────────────────────────────────────────────────────────┘\n\n"

	s += "┌─ MODO ANALYTICS ───────────────────────────────────────────────────┐\n"
	s += "│                                                                      │\n"
	s += "│  Mostra estatísticas completas do grafo atual:                      │\n"
	s += "│  - Total de nós e ligações                                          │\n"
	s += "│  - Distribuição por tipo (Empresas/Pessoas)                         │\n"
	s += "│  - Densidade e grau médio da rede                                   │\n"
	s += "│  - Nós mais conectados (top 10)                                     │\n"
	s += "│  - Tipos de relacionamento                                          │\n"
	s += "│                                                                      │\n"
	s += "│  Q / BACKSPACE  - Voltar ao modo árvore                             │\n"
	s += "│                                                                      │\n"
	s += "└──────────────────────────────────────────────────────────────────────┘\n\n"

	s += "┌─ MODO EXPORTAR ────────────────────────────────────────────────────┐\n"
	s += "│                                                                      │\n"
	s += "│  Exporta o grafo atual em diferentes formatos:                      │\n"
	s += "│                                                                      │\n"
	s += "│  1. Excel (XLSX)     - Arquivo completo com 3 planilhas:            │\n"
	s += "│                        • Nós (entidades)                             │\n"
	s += "│                        • Arestas (relacionamentos)                   │\n"
	s += "│                        • Estatísticas                                │\n"
	s += "│                                                                      │\n"
	s += "│  2. CSV - Nós        - Lista de todas as entidades                  │\n"
	s += "│  3. CSV - Arestas    - Lista de todos os relacionamentos            │\n"
	s += "│  4. CSV - Stats      - Resumo estatístico                           │\n"
	s += "│                                                                      │\n"
	s += "│  ↑↓ Navegar | ENTER Exportar | Q Cancelar                           │\n"
	s += "│  Arquivos salvos em: ./output/                                      │\n"
	s += "│                                                                      │\n"
	s += "└──────────────────────────────────────────────────────────────────────┘\n\n"

	s += "┌─ DICAS ────────────────────────────────────────────────────────────┐\n"
	s += "│                                                                      │\n"
	s += "│  • Comece expandindo o nó raiz para ver os relacionamentos          │\n"
	s += "│  • Cada nó pode ser expandido individualmente                       │\n"
	s += "│  • Use Analytics para ver estatísticas antes de exportar            │\n"
	s += "│  • Ícones: 🏢 = Empresa (PJ), 👤 = Pessoa (PF/PE)                   │\n"
	s += "│  • ▶ = Pode expandir, ▼ = Já expandido                              │\n"
	s += "│                                                                      │\n"
	s += "└──────────────────────────────────────────────────────────────────────┘\n\n"

	s += "Pressione [Q], [BACKSPACE], [F1] ou [?] para voltar\n"

	return s
}
