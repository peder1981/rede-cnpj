package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/graph"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/models"
)

// ServeCaminhos encontra caminhos entre duas entidades
func (h *Handler) ServeCaminhos(c *gin.Context) {
	var req struct {
		From     string `json:"from"`
		To       string `json:"to"`
		MaxDepth int    `json:"maxDepth"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if req.MaxDepth == 0 {
		req.MaxDepth = 5
	}

	pathFinder := graph.NewPathFinder("bases/rede.db")
	result, err := pathFinder.FindPaths(req.From, req.To, req.MaxDepth)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ServeEntidadesComuns encontra entidades em comum
func (h *Handler) ServeEntidadesComuns(c *gin.Context) {
	var req struct {
		ID1 string `json:"id1"`
		ID2 string `json:"id2"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	pathFinder := graph.NewPathFinder("bases/rede.db")
	result, err := pathFinder.FindCommonEntities(req.ID1, req.ID2)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ServeFiltrarGrafo filtra um grafo
func (h *Handler) ServeFiltrarGrafo(c *gin.Context) {
	var req struct {
		Graph    models.Graph `json:"graph"`
		Criteria struct {
			MinConnections int      `json:"minConnections"`
			MaxConnections int      `json:"maxConnections"`
			NodeTypes      []string `json:"nodeTypes"`
			EdgeTypes      []string `json:"edgeTypes"`
		} `json:"criteria"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	criteria := graph.FilterCriteria{
		MinConnections: req.Criteria.MinConnections,
		MaxConnections: req.Criteria.MaxConnections,
		NodeTypes:      req.Criteria.NodeTypes,
		EdgeTypes:      req.Criteria.EdgeTypes,
	}

	filtered := graph.FilterGraph(&req.Graph, criteria)
	c.JSON(http.StatusOK, filtered)
}
