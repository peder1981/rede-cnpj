package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/crossdata"
)

// ServeCrossDataEmpresasPorCPF retorna todas as empresas de um CPF
func (h *Handler) ServeCrossDataEmpresasPorCPF(c *gin.Context) {
	cpf := c.Param("cpf")
	
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	results, err := engine.EmpresasPorCPF(cpf)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"cpf":     cpf,
		"total":   len(results),
		"empresas": results,
	})
}

// ServeCrossDataSociosPorCNPJ retorna todos os sócios de um CNPJ
func (h *Handler) ServeCrossDataSociosPorCNPJ(c *gin.Context) {
	cnpj := c.Param("cnpj")
	
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	results, err := engine.SociosPorCNPJ(cnpj)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"cnpj":   cnpj,
		"total":  len(results),
		"socios": results,
	})
}

// ServeCrossDataSociosEmComum retorna sócios em comum entre duas empresas
func (h *Handler) ServeCrossDataSociosEmComum(c *gin.Context) {
	var req struct {
		CNPJ1 string `json:"cnpj1"`
		CNPJ2 string `json:"cnpj2"`
	}
	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	results, err := engine.SociosEmComum(req.CNPJ1, req.CNPJ2)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"cnpj1":         req.CNPJ1,
		"cnpj2":         req.CNPJ2,
		"total":         len(results),
		"socios_comuns": results,
	})
}

// ServeCrossDataRedeEmpresasPessoa retorna rede de empresas de uma pessoa
func (h *Handler) ServeCrossDataRedeEmpresasPessoa(c *gin.Context) {
	cpf := c.Param("cpf")
	
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	results, err := engine.RedeEmpresasPessoa(cpf)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"cpf":   cpf,
		"total": len(results),
		"rede":  results,
	})
}

// ServeCrossDataEmpresasMesmoEndereco retorna empresas no mesmo endereço
func (h *Handler) ServeCrossDataEmpresasMesmoEndereco(c *gin.Context) {
	var req struct {
		CEP        string `json:"cep"`
		Logradouro string `json:"logradouro"`
		Numero     string `json:"numero"`
	}
	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	results, err := engine.EmpresasMesmoEndereco(req.CEP, req.Logradouro, req.Numero)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"endereco": req,
		"total":    len(results),
		"empresas": results,
	})
}

// ServeCrossDataEmpresasMesmoContato retorna empresas com mesmo email/telefone
func (h *Handler) ServeCrossDataEmpresasMesmoContato(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Telefone string `json:"telefone"`
	}
	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	results, err := engine.EmpresasMesmoContato(req.Email, req.Telefone)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"filtro":   req,
		"total":    len(results),
		"empresas": results,
	})
}

// ServeCrossDataRepresentantesLegais retorna menores com representantes
func (h *Handler) ServeCrossDataRepresentantesLegais(c *gin.Context) {
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	results, err := engine.RepresentantesLegais()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"total":           len(results),
		"representantes": results,
	})
}

// ServeCrossDataEmpresasEstrangeiras retorna empresas estrangeiras
func (h *Handler) ServeCrossDataEmpresasEstrangeiras(c *gin.Context) {
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	results, err := engine.EmpresasEstrangeiras()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"total":    len(results),
		"empresas": results,
	})
}

// ServeCrossDataSociosEstrangeiros retorna sócios estrangeiros
func (h *Handler) ServeCrossDataSociosEstrangeiros(c *gin.Context) {
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	results, err := engine.SociosEstrangeiros()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"total":  len(results),
		"socios": results,
	})
}

// ServeCrossDataTimelinePessoa retorna timeline de atividades de uma pessoa
func (h *Handler) ServeCrossDataTimelinePessoa(c *gin.Context) {
	cpf := c.Param("cpf")
	
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	results, err := engine.TimelinePessoa(cpf)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"cpf":      cpf,
		"total":    len(results),
		"timeline": results,
	})
}

// ServeCrossDataSociosEmpresasBaixadas retorna sócios com empresas baixadas
func (h *Handler) ServeCrossDataSociosEmpresasBaixadas(c *gin.Context) {
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	results, err := engine.SociosEmpresasBaixadas()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"total":  len(results),
		"socios": results,
	})
}

// ServeCrossDataDadosCompletos retorna dados completos de empresa SEM CENSURA
func (h *Handler) ServeCrossDataDadosCompletos(c *gin.Context) {
	cnpj := c.Param("cnpj")
	
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	result, err := engine.DadosCompletosEmpresa(cnpj)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, result)
}
