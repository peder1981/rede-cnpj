package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/config"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/database"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/services"
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

	printHeader()
	fmt.Printf("📊 Banco de Dados: %s\n", cfg.ReferenciaBD)
	fmt.Println("")

	// Cria serviço
	redeService := services.NewRedeService(cfg)
	
	// Pergunta o CNPJ inicial
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("🔍 Digite o CNPJ/CPF inicial: ")
	
	if !scanner.Scan() {
		return
	}
	
	cnpj := strings.TrimSpace(scanner.Text())
	if cnpj == "" {
		fmt.Println("❌ CNPJ/CPF não pode ser vazio")
		return
	}

	// Inicia interface TUI
	fmt.Println("\n⏳ Carregando interface interativa...")
	p := tea.NewProgram(initialModel(redeService, cnpj), tea.WithAltScreen())
	
	if _, err := p.Run(); err != nil {
		log.Fatalf("Erro ao executar TUI: %v", err)
	}
}

func printHeader() {
	fmt.Println("")
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                                ║")
	fmt.Println("║         🔍  RedeCNPJ - Navegação Interativa                   ║")
	fmt.Println("║                                                                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println("")
}
