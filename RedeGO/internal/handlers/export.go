package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/export"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/models"
)

// ServeExportExcel exporta grafo para Excel
func (h *Handler) ServeExportExcel(c *gin.Context) {
	var graph models.Graph
	if err := c.BindJSON(&graph); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	exporter := export.NewExcelExporter()
	defer exporter.Close()

	data, err := exporter.ExportGraph(&graph)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=rede-cnpj.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

// ServeExportCSV exporta grafo para CSV
func (h *Handler) ServeExportCSV(c *gin.Context) {
	var graph models.Graph
	if err := c.BindJSON(&graph); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	exporter := export.NewCSVExporter()
	
	// Exporta apenas nós ou arestas dependendo do parâmetro
	tipo := c.Query("tipo")
	
	var data []byte
	var err error
	var filename string

	switch tipo {
	case "nos":
		data, err = exporter.ExportNodes(graph.Nodes)
		filename = "nos.csv"
	case "arestas":
		data, err = exporter.ExportEdges(graph.Edges)
		filename = "arestas.csv"
	case "stats":
		data, err = exporter.ExportStats(&graph)
		filename = "estatisticas.csv"
	default:
		// Por padrão, exporta nós
		data, err = exporter.ExportNodes(graph.Nodes)
		filename = "nos.csv"
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/csv", data)
}
