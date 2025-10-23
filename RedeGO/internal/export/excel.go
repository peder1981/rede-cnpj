package export

import (
	"fmt"

	"github.com/peder1981/rede-cnpj/RedeGO/internal/models"
	"github.com/xuri/excelize/v2"
)

// ExcelExporter exporta dados para Excel
type ExcelExporter struct {
	file *excelize.File
}

// NewExcelExporter cria um novo exportador Excel
func NewExcelExporter() *ExcelExporter {
	return &ExcelExporter{
		file: excelize.NewFile(),
	}
}

// ExportGraph exporta grafo para Excel
func (e *ExcelExporter) ExportGraph(graph *models.Graph) ([]byte, error) {
	// Cria planilha de nós
	if err := e.createNodesSheet(graph.Nodes); err != nil {
		return nil, err
	}

	// Cria planilha de arestas
	if err := e.createEdgesSheet(graph.Edges); err != nil {
		return nil, err
	}

	// Cria planilha de estatísticas
	if err := e.createStatsSheet(graph); err != nil {
		return nil, err
	}

	// Salva em buffer
	buf, err := e.file.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// createNodesSheet cria planilha de nós
func (e *ExcelExporter) createNodesSheet(nodes []models.Node) error {
	sheetName := "Nós"
	index, err := e.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Cabeçalhos
	headers := []string{"ID", "Label", "Tipo"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		e.file.SetCellValue(sheetName, cell, h)
	}

	// Estilo do cabeçalho
	style, _ := e.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4472C4"}, Pattern: 1},
	})
	e.file.SetCellStyle(sheetName, "A1", "C1", style)

	// Dados
	for i, node := range nodes {
		row := i + 2
		tipo := "Empresa"
		if len(node.ID) > 3 && (node.ID[:3] == "PF_" || node.ID[:3] == "PE_") {
			tipo = "Pessoa"
		}

		e.file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), node.ID)
		e.file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), node.Label)
		e.file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), tipo)
	}

	// Ajusta largura das colunas
	e.file.SetColWidth(sheetName, "A", "A", 30)
	e.file.SetColWidth(sheetName, "B", "B", 50)
	e.file.SetColWidth(sheetName, "C", "C", 15)

	e.file.SetActiveSheet(index)
	return nil
}

// createEdgesSheet cria planilha de arestas
func (e *ExcelExporter) createEdgesSheet(edges []models.Edge) error {
	sheetName := "Arestas"
	_, err := e.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Cabeçalhos
	headers := []string{"Origem", "Destino", "Tipo"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		e.file.SetCellValue(sheetName, cell, h)
	}

	// Estilo do cabeçalho
	style, _ := e.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#70AD47"}, Pattern: 1},
	})
	e.file.SetCellStyle(sheetName, "A1", "C1", style)

	// Dados
	for i, edge := range edges {
		row := i + 2
		e.file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), edge.From)
		e.file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), edge.To)
		e.file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), edge.Label)
	}

	// Ajusta largura das colunas
	e.file.SetColWidth(sheetName, "A", "B", 30)
	e.file.SetColWidth(sheetName, "C", "C", 25)

	return nil
}

// createStatsSheet cria planilha de estatísticas
func (e *ExcelExporter) createStatsSheet(graph *models.Graph) error {
	sheetName := "Estatísticas"
	_, err := e.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Conta tipos
	empresas := 0
	pessoas := 0
	for _, node := range graph.Nodes {
		if len(node.ID) > 3 && (node.ID[:3] == "PF_" || node.ID[:3] == "PE_") {
			pessoas++
		} else {
			empresas++
		}
	}

	// Dados
	stats := [][]interface{}{
		{"Métrica", "Valor"},
		{"Total de Nós", len(graph.Nodes)},
		{"Total de Arestas", len(graph.Edges)},
		{"Empresas (PJ)", empresas},
		{"Pessoas (PF/PE)", pessoas},
	}

	for i, row := range stats {
		for j, val := range row {
			cell := fmt.Sprintf("%c%d", 'A'+j, i+1)
			e.file.SetCellValue(sheetName, cell, val)
		}
	}

	// Estilo do cabeçalho
	style, _ := e.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFC000"}, Pattern: 1},
	})
	e.file.SetCellStyle(sheetName, "A1", "B1", style)

	// Ajusta largura das colunas
	e.file.SetColWidth(sheetName, "A", "A", 20)
	e.file.SetColWidth(sheetName, "B", "B", 15)

	// Remove planilha padrão
	e.file.DeleteSheet("Sheet1")

	return nil
}

// Close fecha o arquivo
func (e *ExcelExporter) Close() error {
	return e.file.Close()
}
