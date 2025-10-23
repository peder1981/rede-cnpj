package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/search"
)

// ServeBuscaAvancada busca avançada com FTS5
func (h *Handler) ServeBuscaAvancada(c *gin.Context) {
	var req struct {
		Query      string `json:"query"`
		Limit      int    `json:"limit"`
		UseGlob    bool   `json:"useGlob"`
		RandomTest bool   `json:"randomTest"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// Limite padrão
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Caminho do banco de busca
	searchDB := "bases/rede_search.db"

	// Cria buscador
	searcher := search.NewAdvancedSearch(searchDB)

	// Executa busca
	results, err := searcher.Search(search.SearchOptions{
		Query:      req.Query,
		Limit:      req.Limit,
		UseGlob:    req.UseGlob,
		RandomTest: req.RandomTest,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"count":   len(results),
	})
}
