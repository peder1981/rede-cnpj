package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/config"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/models"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/services"
)

// Handler gerencia as requisições HTTP
type Handler struct {
	cfg         *config.Config
	redeService *services.RedeService
}

// NewHandler cria uma nova instância do handler
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg:         cfg,
		redeService: services.NewRedeService(cfg),
	}
}

// ServeHTMLPagina renderiza a página principal
func (h *Handler) ServeHTMLPagina(c *gin.Context) {
	cpfcnpj := c.Param("cpfcnpj")
	camadaStr := c.Param("camada")
	idArquivoServidor := c.Param("idArquivoServidor")

	camada := 0
	if camadaStr != "" {
		camada, _ = strconv.Atoi(camadaStr)
	}
	if camada > 10 {
		camada = 10
	}

	params := gin.H{
		"cpfcnpj":              cpfcnpj,
		"camada":               camada,
		"idArquivoServidor":    idArquivoServidor,
		"bBaseReceita":         h.cfg.BaseReceita != "",
		"bBaseLocal":           h.cfg.BaseLocal != "",
		"btextoEmbaixoIcone":   h.cfg.BTextoEmbaixoIcone,
		"referenciaBD":         h.cfg.ReferenciaBD,
		"geocode_max":          h.cfg.GeocodeMax,
		"bbusca_chaves":        h.cfg.BuscaChaves,
		"usuarioLocal":         isLocalUser(c),
		"mobile":               isMobile(c),
		"chrome":               strings.Contains(c.Request.UserAgent(), "Chrome"),
		"firefox":              strings.Contains(c.Request.UserAgent(), "Firefox"),
		"bMenuInserirInicial":  h.cfg.ExibeMenuInserir,
		"mensagem":             "",
	}

	c.HTML(http.StatusOK, "rede_template.html", params)
}

// ServeRedeJSONCNPJ retorna dados da rede em formato JSON
func (h *Handler) ServeRedeJSONCNPJ(c *gin.Context) {
	tipo := c.Param("tipo")
	camadaStr := c.Param("camada")
	cpfcnpj := c.Param("cpfcnpj")

	camada, _ := strconv.Atoi(camadaStr)
	if camada > 10 {
		camada = 10
	}

	var listaIDs []string
	if c.Request.Method == "POST" {
		if err := c.BindJSON(&listaIDs); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}
	} else {
		listaIDs = []string{cpfcnpj}
	}

	criterioCaminhos := ""
	if strings.HasPrefix(tipo, "caminhos") {
		criterioCaminhos = strings.TrimPrefix(tipo, "caminhos-")
	}

	graph, err := h.redeService.CamadasRede(camada, listaIDs, "", criterioCaminhos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, graph)
}

// ServeDadosDetalhes retorna dados detalhados de um CNPJ
func (h *Handler) ServeDadosDetalhes(c *gin.Context) {
	cpfcnpj := c.Param("cpfcnpj")

	if c.Request.Method == "POST" {
		var req struct {
			IDIn string `json:"idin"`
		}
		if err := c.BindJSON(&req); err == nil {
			cpfcnpj = req.IDIn
		}
	}

	dados := h.redeService.GetDadosCNPJ(cpfcnpj)
	if dados == nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	c.JSON(http.StatusOK, dados)
}

// ServeBuscaPorNome busca empresas ou sócios por nome
func (h *Handler) ServeBuscaPorNome(c *gin.Context) {
	nome := c.Query("q")
	limiteStr := c.Query("limite")

	limite := 10
	if limiteStr != "" {
		limite, _ = strconv.Atoi(limiteStr)
	}

	if nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'q' é obrigatório"})
		return
	}

	results, err := h.redeService.BuscaPorNome(nome, limite)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// ServeArquivosJSON serve arquivos JSON salvos
func (h *Handler) ServeArquivosJSON(c *gin.Context) {
	arquivopath := c.Param("arquivopath")
	
	// Implementação simplificada - na versão completa, adicionar lógica de segurança
	c.File(fmt.Sprintf("%s/%s", h.cfg.PastaArquivos, arquivopath))
}

// ServeArquivosJSONUpload faz upload de arquivos JSON
func (h *Handler) ServeArquivosJSONUpload(c *gin.Context) {
	nomeArquivo := c.Param("nomeArquivo")

	var data interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, models.FileUploadResponse{
			Mensagem: "JSON inválido",
		})
		return
	}

	// Verifica tamanho
	jsonBytes, _ := json.Marshal(data)
	if len(jsonBytes) > 100000 {
		c.JSON(http.StatusOK, models.FileUploadResponse{
			Mensagem: "O arquivo é muito grande e não foi salvo",
		})
		return
	}

	// Salva arquivo (implementação simplificada)
	c.JSON(http.StatusOK, models.FileUploadResponse{
		NomeArquivoServidor: nomeArquivo,
	})
}

// ServeDadosEmArquivo exporta dados para Excel ou outros formatos
func (h *Handler) ServeDadosEmArquivo(c *gin.Context) {
	formato := c.Param("formato")

	var dados models.ExportRequest
	if err := c.Bind(&dados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	switch formato {
	case "xlsx":
		// Implementar exportação para Excel
		c.JSON(http.StatusOK, gin.H{"message": "Exportação Excel não implementada ainda"})
	case "anx":
		// Implementar exportação para i2
		c.JSON(http.StatusOK, gin.H{"message": "Exportação i2 não implementada ainda"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato não suportado"})
	}
}

// ServeMapa gera mapa com endereços
func (h *Handler) ServeMapa(c *gin.Context) {
	var req models.MapaRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Implementação simplificada
	c.JSON(http.StatusOK, gin.H{"message": "Geração de mapa não implementada ainda"})
}

// ServeDadosPublicosDisponivel verifica dados públicos disponíveis
func (h *Handler) ServeDadosPublicosDisponivel(c *gin.Context) {
	// Implementação simplificada - na versão completa, fazer scraping do site da Receita
	c.JSON(http.StatusOK, models.DadosPublicosResponse{
		AnoMesSendoUsado: "2024-10",
		AnoMesDisponivel: "2024-10",
		URL:              "https://arquivos.receitafederal.gov.br/cnpj/dados_abertos_cnpj/",
	})
}

// ServeAPIStatus retorna status da API
func (h *Handler) ServeAPIStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
		"message": "RedeCNPJ API em Go",
	})
}

// Funções auxiliares

func isLocalUser(c *gin.Context) bool {
	return c.ClientIP() == "127.0.0.1" || c.ClientIP() == "::1"
}

func isMobile(c *gin.Context) bool {
	ua := c.Request.UserAgent()
	return strings.Contains(ua, "Mobile") || 
	       strings.Contains(ua, "Opera Mini") || 
	       strings.Contains(ua, "Android")
}
