package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/peder1981/rede-cnpj/RedeGO/internal/models"
	"github.com/tealeg/xlsx/v3"
)

// ExportService gerencia exportação de dados
type ExportService struct{}

// NewExportService cria uma nova instância do serviço de exportação
func NewExportService() *ExportService {
	return &ExportService{}
}

// ExportToExcel exporta dados para formato Excel
func (s *ExportService) ExportToExcel(graph *models.Graph) (*bytes.Buffer, error) {
	file := xlsx.NewFile()

	// Planilha de Nós
	sheetNodes, err := file.AddSheet("Nós")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar planilha de nós: %w", err)
	}

	// Cabeçalho
	headerRow := sheetNodes.AddRow()
	headers := []string{"ID", "Label", "Tipo", "Cor", "Ícone", "Nota", "Camada"}
	for _, h := range headers {
		cell := headerRow.AddCell()
		cell.Value = h
		cell.GetStyle().Font.Bold = true
	}

	// Dados dos nós
	for _, node := range graph.Nodes {
		row := sheetNodes.AddRow()
		row.AddCell().Value = node.ID
		row.AddCell().Value = node.Label
		row.AddCell().Value = node.Type
		row.AddCell().Value = node.Color
		row.AddCell().Value = node.Icon
		row.AddCell().Value = node.Note
		row.AddCell().SetInt(node.Camada)
	}

	// Planilha de Ligações
	sheetEdges, err := file.AddSheet("Ligações")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar planilha de ligações: %w", err)
	}

	// Cabeçalho
	headerRow = sheetEdges.AddRow()
	headers = []string{"De", "Para", "Label", "Tipo", "Qualificação", "Valor"}
	for _, h := range headers {
		cell := headerRow.AddCell()
		cell.Value = h
		cell.GetStyle().Font.Bold = true
	}

	// Dados das ligações
	for _, edge := range graph.Edges {
		row := sheetEdges.AddRow()
		row.AddCell().Value = edge.From
		row.AddCell().Value = edge.To
		row.AddCell().Value = edge.Label
		row.AddCell().Value = edge.Type
		row.AddCell().Value = edge.Qualificacao
		row.AddCell().SetFloat(edge.Value)
	}

	// Escreve para buffer
	var buffer bytes.Buffer
	if err := file.Write(&buffer); err != nil {
		return nil, fmt.Errorf("erro ao escrever Excel: %w", err)
	}

	return &buffer, nil
}

// ExportToJSON exporta dados para formato JSON
func (s *ExportService) ExportToJSON(graph *models.Graph) ([]byte, error) {
	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar JSON: %w", err)
	}
	return data, nil
}

// ExportToCSV exporta dados para formato CSV
func (s *ExportService) ExportToCSV(graph *models.Graph) ([]byte, error) {
	var buffer bytes.Buffer

	// CSV de Nós
	buffer.WriteString("# Nós\n")
	buffer.WriteString("ID,Label,Tipo,Cor,Ícone,Nota,Camada\n")
	for _, node := range graph.Nodes {
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%d\n",
			escapeCsv(node.ID),
			escapeCsv(node.Label),
			escapeCsv(node.Type),
			escapeCsv(node.Color),
			escapeCsv(node.Icon),
			escapeCsv(node.Note),
			node.Camada,
		)
		buffer.WriteString(line)
	}

	buffer.WriteString("\n# Ligações\n")
	buffer.WriteString("De,Para,Label,Tipo,Qualificação,Valor\n")
	for _, edge := range graph.Edges {
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%.2f\n",
			escapeCsv(edge.From),
			escapeCsv(edge.To),
			escapeCsv(edge.Label),
			escapeCsv(edge.Type),
			escapeCsv(edge.Qualificacao),
			edge.Value,
		)
		buffer.WriteString(line)
	}

	return buffer.Bytes(), nil
}

// ExportToI2 exporta dados para formato i2 Chart Reader (.anx)
func (s *ExportService) ExportToI2(graph *models.Graph) (*bytes.Buffer, error) {
	// Estrutura simplificada do formato i2
	// Na versão completa, implementar XML completo do formato ANX
	
	i2Data := map[string]interface{}{
		"chart": map[string]interface{}{
			"version": "1.0",
			"created": time.Now().Format(time.RFC3339),
			"entities": s.convertNodesToI2Entities(graph.Nodes),
			"links":    s.convertEdgesToI2Links(graph.Edges),
		},
	}

	data, err := json.MarshalIndent(i2Data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar i2: %w", err)
	}

	return bytes.NewBuffer(data), nil
}

// convertNodesToI2Entities converte nós para formato i2
func (s *ExportService) convertNodesToI2Entities(nodes []models.Node) []map[string]interface{} {
	entities := make([]map[string]interface{}, len(nodes))
	
	for i, node := range nodes {
		entityType := "Person"
		if node.Type == "PJ" {
			entityType = "Organization"
		}

		entities[i] = map[string]interface{}{
			"id":    node.ID,
			"label": node.Label,
			"type":  entityType,
			"properties": map[string]interface{}{
				"icon":  node.Icon,
				"color": node.Color,
				"note":  node.Note,
			},
		}
	}

	return entities
}

// convertEdgesToI2Links converte arestas para formato i2
func (s *ExportService) convertEdgesToI2Links(edges []models.Edge) []map[string]interface{} {
	links := make([]map[string]interface{}, len(edges))
	
	for i, edge := range edges {
		links[i] = map[string]interface{}{
			"from":  edge.From,
			"to":    edge.To,
			"label": edge.Label,
			"type":  edge.Type,
			"properties": map[string]interface{}{
				"qualificacao": edge.Qualificacao,
				"valor":        edge.Value,
			},
		}
	}

	return links
}

// escapeCsv escapa valores para CSV
func escapeCsv(s string) string {
	// Se contém vírgula, aspas ou quebra de linha, envolve em aspas
	needsQuotes := false
	for _, c := range s {
		if c == ',' || c == '"' || c == '\n' || c == '\r' {
			needsQuotes = true
			break
		}
	}

	if !needsQuotes {
		return s
	}

	// Duplica aspas internas
	result := ""
	for _, c := range s {
		if c == '"' {
			result += "\"\""
		} else {
			result += string(c)
		}
	}

	return "\"" + result + "\""
}

// GetExportStats retorna estatísticas do grafo
func (s *ExportService) GetExportStats(graph *models.Graph) map[string]interface{} {
	stats := map[string]interface{}{
		"total_nodes": len(graph.Nodes),
		"total_edges": len(graph.Edges),
		"node_types":  make(map[string]int),
		"edge_types":  make(map[string]int),
	}

	nodeTypes := make(map[string]int)
	for _, node := range graph.Nodes {
		nodeTypes[node.Type]++
	}
	stats["node_types"] = nodeTypes

	edgeTypes := make(map[string]int)
	for _, edge := range graph.Edges {
		edgeTypes[edge.Type]++
	}
	stats["edge_types"] = edgeTypes

	return stats
}
