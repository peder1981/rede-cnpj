package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/peder1981/rede-cnpj/RedeGO/internal/config"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/importer"
)

func main() {
	// Flags
	downloadOnly := flag.Bool("download", false, "Apenas baixa os arquivos ZIP")
	processOnly := flag.Bool("process", false, "Apenas processa arquivos já baixados")
	createLinks := flag.Bool("links", false, "Cria tabelas de ligação (rede.db)")
	createSearch := flag.Bool("search", false, "Cria índices de busca (rede_search.db)")
	all := flag.Bool("all", false, "Executa todo o processo (download + process + links + search)")
	confFile := flag.String("config", "rede.ini", "Arquivo de configuração (opcional)")
	
	flag.Parse()

	// Tenta carregar configuração do arquivo
	var cfg *config.Config
	var err error
	
	// Verifica se arquivo de config existe
	if _, statErr := os.Stat(*confFile); statErr == nil {
		fmt.Printf("ℹ️  Carregando configuração de %s...\n", *confFile)
		cfg, err = config.LoadConfig()
		if err != nil {
			fmt.Printf("⚠️  Erro ao carregar config: %v\n", err)
			fmt.Println("⚠️  Usando configuração padrão (SQLite)")
			cfg = &config.Config{
				BaseReceita:    "bases/cnpj.db",
				BaseRede:       "bases/rede.db",
				BaseRedeSearch: "bases/rede_search.db",
				PastaArquivos:  "arquivos",
			}
		}
	} else {
		// Arquivo não existe, usa config padrão
		fmt.Println("ℹ️  Nenhum arquivo de configuração encontrado, usando SQLite")
		cfg = &config.Config{
			BaseReceita:    "bases/cnpj.db",
			BaseRede:       "bases/rede.db",
			BaseRedeSearch: "bases/rede_search.db",
			PastaArquivos:  "arquivos",
		}
	}
	
	// Exibe tipo de banco que será usado
	if cfg.PostgresURL != "" {
		fmt.Println("🐘 Banco de dados: PostgreSQL (importação direta)")
	} else {
		fmt.Println("🗄️  Banco de dados: SQLite (modo legado)")
	}

	// Cria importador
	imp := importer.NewImporter(cfg)

	printHeader()

	// Determina o que executar
	if *all {
		*downloadOnly = false
		*processOnly = false
		*createLinks = true
		*createSearch = true
		runAll(imp)
	} else if *downloadOnly {
		if err := imp.DownloadFiles(); err != nil {
			log.Fatalf("Erro no download: %v", err)
		}
	} else if *processOnly {
		if err := imp.ProcessFiles(); err != nil {
			log.Fatalf("Erro no processamento: %v", err)
		}
	} else if *createLinks {
		if err := imp.CreateLinkTables(); err != nil {
			log.Fatalf("Erro ao criar tabelas de ligação: %v", err)
		}
	} else if *createSearch {
		if err := imp.CreateSearchIndexes(); err != nil {
			log.Fatalf("Erro ao criar índices de busca: %v", err)
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("\n✅ Processo concluído com sucesso!")
}

func runAll(imp *importer.Importer) {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Download dos arquivos", imp.DownloadFiles},
		{"Processamento dos arquivos", imp.ProcessFiles},
		{"Criação de tabelas de ligação", imp.CreateLinkTables},
		{"Criação de índices de busca", imp.CreateSearchIndexes},
	}

	for i, step := range steps {
		fmt.Printf("\n[%d/%d] %s...\n", i+1, len(steps), step.name)
		if err := step.fn(); err != nil {
			log.Fatalf("❌ Erro em '%s': %v", step.name, err)
		}
		fmt.Printf("✅ %s concluído\n", step.name)
	}
}

func printHeader() {
	fmt.Println("")
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                                ║")
	fmt.Println("║     📥 RedeCNPJ - Importador de Dados da Receita Federal      ║")
	fmt.Println("║                                                                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println("")
}
