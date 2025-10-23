package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/analytics"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/models"
)

// ServeAnalytics retorna estatísticas do grafo
func (h *Handler) ServeAnalytics(c *gin.Context) {
	var graph models.Graph
	if err := c.BindJSON(&graph); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	analyzer := analytics.NewAnalyzer()
	stats := analyzer.AnalyzeGraph(&graph)

	c.JSON(http.StatusOK, stats)
}

// ServeNosCentrais retorna nós centrais
func (h *Handler) ServeNosCentrais(c *gin.Context) {
	var req struct {
		Graph models.Graph `json:"graph"`
		Top   int          `json:"top"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if req.Top == 0 {
		req.Top = 10
	}

	analyzer := analytics.NewAnalyzer()
	central := analyzer.DetectCentralNodes(&req.Graph, req.Top)

	c.JSON(http.StatusOK, gin.H{"centralNodes": central})
}

// ServeComunidades detecta comunidades
func (h *Handler) ServeComunidades(c *gin.Context) {
	var graph models.Graph
	if err := c.BindJSON(&graph); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	analyzer := analytics.NewAnalyzer()
	communities := analyzer.DetectCommunities(&graph)

	c.JSON(http.StatusOK, gin.H{"communities": communities})
}

// ServeCaminhoMaisCurto encontra caminho mais curto
func (h *Handler) ServeCaminhoMaisCurto(c *gin.Context) {
	var req struct {
		Graph models.Graph `json:"graph"`
		From  string       `json:"from"`
		To    string       `json:"to"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	analyzer := analytics.NewAnalyzer()
	path := analyzer.CalculateShortestPath(&req.Graph, req.From, req.To)

	c.JSON(http.StatusOK, gin.H{"path": path})
}
