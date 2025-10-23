package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/forensics"
)

// ServeForensicsInvestigatePerson perfil completo de suspeito
func (h *Handler) ServeForensicsInvestigatePerson(c *gin.Context) {
	cpf := c.Param("cpf")
	
	inv := forensics.NewInvestigator("bases/cnpj.db", "bases/rede.db")
	profile, err := inv.InvestigatePerson(cpf)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, profile)
}

// ServeForensicsShellCompanies detecta empresas de fachada
func (h *Handler) ServeForensicsShellCompanies(c *gin.Context) {
	minStr := c.DefaultQuery("min_empresas", "10")
	minEmpresas, _ := strconv.Atoi(minStr)
	
	inv := forensics.NewInvestigator("bases/cnpj.db", "bases/rede.db")
	clusters, err := inv.DetectShellCompanies(minEmpresas)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"total":    len(clusters),
		"clusters": clusters,
	})
}

// ServeForensicsFrontmen detecta laranjas
func (h *Handler) ServeForensicsFrontmen(c *gin.Context) {
	var req struct {
		Criterio string `json:"criterio"` // "telefone" ou "email"
		Valor    string `json:"valor"`
	}
	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	
	inv := forensics.NewInvestigator("bases/cnpj.db", "bases/rede.db")
	cluster, err := inv.DetectFrontmen(req.Criterio, req.Valor)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, cluster)
}

// ServeForensicsMassRegistration detecta abertura em massa
func (h *Handler) ServeForensicsMassRegistration(c *gin.Context) {
	cpf := c.Param("cpf")
	diasStr := c.DefaultQuery("dias", "30")
	dias, _ := strconv.Atoi(diasStr)
	
	inv := forensics.NewInvestigator("bases/cnpj.db", "bases/rede.db")
	results, err := inv.DetectMassRegistration(cpf, dias)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"cpf":     cpf,
		"total":   len(results),
		"eventos": results,
	})
}

// ServeForensicsOwnershipChain rastreia cadeia de controle
func (h *Handler) ServeForensicsOwnershipChain(c *gin.Context) {
	cnpj := c.Param("cnpj")
	nivelStr := c.DefaultQuery("max_nivel", "3")
	maxNivel, _ := strconv.Atoi(nivelStr)
	
	inv := forensics.NewInvestigator("bases/cnpj.db", "bases/rede.db")
	chain, err := inv.TraceOwnershipChain(cnpj, maxNivel)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"cnpj":    cnpj,
		"niveis":  maxNivel,
		"total":   len(chain),
		"cadeia":  chain,
	})
}

// ServeForensicsSuspiciousPatterns detecta padrões suspeitos
func (h *Handler) ServeForensicsSuspiciousPatterns(c *gin.Context) {
	inv := forensics.NewInvestigator("bases/cnpj.db", "bases/rede.db")
	patterns, err := inv.DetectSuspiciousPatterns()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"total":    len(patterns),
		"suspeitos": patterns,
	})
}
