package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/config"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/database"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/handlers"
)

func main() {
	// Carrega configuração
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	// Inicializa bancos de dados
	if err := database.InitDatabases(cfg); err != nil {
		log.Fatalf("Erro ao inicializar bancos de dados: %v", err)
	}
	defer database.Close()

	// Carrega dicionários
	if _, err := database.LoadDicionarios(); err != nil {
		log.Printf("AVISO: Erro ao carregar dicionários: %v", err)
	}

	// Configura Gin
	if !cfg.BuscaGoogle {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Carrega templates HTML
	router.LoadHTMLGlob("templates/*")

	// Serve arquivos estáticos
	router.Static("/rede/static", "./static")
	router.Static("/rede/arquivos", cfg.PastaArquivos)

	// Cria handlers
	h := handlers.NewHandler(cfg)

	// Rotas principais
	router.GET("/rede/", h.ServeHTMLPagina)
	router.POST("/rede/", h.ServeHTMLPagina)
	router.GET("/rede/grafico/:camada/:cpfcnpj", h.ServeHTMLPagina)
	router.GET("/rede/grafico_no_servidor/:idArquivoServidor", h.ServeHTMLPagina)

	// API de dados
	router.POST("/rede/grafojson/:tipo/:camada/:cpfcnpj", h.ServeRedeJSONCNPJ)
	router.GET("/rede/dadosjson/:cpfcnpj", h.ServeDadosDetalhes)
	router.POST("/rede/dadosjson/:cpfcnpj", h.ServeDadosDetalhes)

	// Busca
	router.GET("/rede/busca", h.ServeBuscaPorNome)

	// Arquivos
	router.GET("/rede/arquivos_json/:arquivopath", h.ServeArquivosJSON)
	router.POST("/rede/arquivos_json/:arquivopath", h.ServeArquivosJSON)
	router.DELETE("/rede/arquivos_json/:arquivopath", h.ServeArquivosJSON)
	router.POST("/rede/arquivos_json_upload/:nomeArquivo", h.ServeArquivosJSONUpload)

	// Exportação
	router.POST("/rede/dadosemarquivo/:formato", h.ServeDadosEmArquivo)
	router.POST("/rede/mapa", h.ServeMapa)

	// Informações
	router.GET("/rede/informacao/dados_publicos_cnpj_disponivel", h.ServeDadosPublicosDisponivel)
	router.GET("/rede/api/status", h.ServeAPIStatus)

	// Configuração de shutdown gracioso
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Encerrando servidor...")
		database.Close()
		os.Exit(0)
	}()

	// Inicia servidor
	addr := fmt.Sprintf(":%d", cfg.PortaFlask)
	log.Printf("Servidor RedeCNPJ iniciado em http://127.0.0.1%s/rede/", addr)
	log.Printf("Referência BD: %s", cfg.ReferenciaBD)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
