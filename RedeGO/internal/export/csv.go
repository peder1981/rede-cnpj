package export

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/peder1981/rede-cnpj/RedeGO/internal/models"
)

// CSVExporter exporta dados para CSV
type CSVExporter struct{}

// NewCSVExporter cria um novo exportador CSV
func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

// ExportNodes exporta nós para CSV
func (e *CSVExporter) ExportNodes(nodes []models.Node) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Cabeçalho
	if err := w.Write([]string{"ID", "Label", "Tipo"}); err != nil {
		return nil, err
	}

	// Dados
	for _, node := range nodes {
		tipo := "Empresa"
		if len(node.ID) > 3 && (node.ID[:3] == "PF_" || node.ID[:3] == "PE_") {
			tipo = "Pessoa"
		}

		record := []string{node.ID, node.Label, tipo}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ExportEdges exporta arestas para CSV
func (e *CSVExporter) ExportEdges(edges []models.Edge) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Cabeçalho
	if err := w.Write([]string{"Origem", "Destino", "Tipo"}); err != nil {
		return nil, err
	}

	// Dados
	for _, edge := range edges {
		record := []string{edge.From, edge.To, edge.Label}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ExportGraph exporta grafo completo para CSV (múltiplos arquivos em ZIP)
func (e *CSVExporter) ExportGraph(graph *models.Graph) (map[string][]byte, error) {
	files := make(map[string][]byte)

	// Exporta nós
	nodesData, err := e.ExportNodes(graph.Nodes)
	if err != nil {
		return nil, fmt.Errorf("erro ao exportar nós: %w", err)
	}
	files["nos.csv"] = nodesData

	// Exporta arestas
	edgesData, err := e.ExportEdges(graph.Edges)
	if err != nil {
		return nil, fmt.Errorf("erro ao exportar arestas: %w", err)
	}
	files["arestas.csv"] = edgesData

	// Exporta estatísticas
	statsData, err := e.ExportStats(graph)
	if err != nil {
		return nil, fmt.Errorf("erro ao exportar estatísticas: %w", err)
	}
	files["estatisticas.csv"] = statsData

	return files, nil
}

// ExportStats exporta estatísticas para CSV
func (e *CSVExporter) ExportStats(graph *models.Graph) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Cabeçalho
	if err := w.Write([]string{"Métrica", "Valor"}); err != nil {
		return nil, err
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
	stats := [][]string{
		{"Total de Nós", fmt.Sprintf("%d", len(graph.Nodes))},
		{"Total de Arestas", fmt.Sprintf("%d", len(graph.Edges))},
		{"Empresas (PJ)", fmt.Sprintf("%d", empresas)},
		{"Pessoas (PF/PE)", fmt.Sprintf("%d", pessoas)},
	}

	for _, record := range stats {
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
