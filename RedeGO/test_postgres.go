package main

import (
	"fmt"
	"log"

	"github.com/peder1981/rede-cnpj/RedeGO/internal/config"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/database"
)

func main() {
	fmt.Println("🧪 Teste de Conexão PostgreSQL")
	fmt.Println("")

	// Criar config com PostgreSQL
	cfg := &config.Config{
		PostgresURL: "postgresql://rede_user:rede_cnpj_2025@localhost:5433/rede_cnpj?sslmode=disable",
	}

	// Inicializar database
	err := database.InitDatabases(cfg)
	if err != nil {
		log.Fatalf("❌ Erro ao inicializar: %v", err)
	}
	defer database.Close()

	fmt.Println("")
	fmt.Println("📊 Testando queries...")
	fmt.Println("")

	// Testar conexão
	db := database.GetDB()

	// Contar empresas
	var countEmpresas int64
	err = db.QueryRow("SELECT COUNT(*) FROM receita.empresas").Scan(&countEmpresas)
	if err != nil {
		log.Printf("⚠️  Erro ao contar empresas: %v", err)
	} else {
		fmt.Printf("✅ Empresas: %d\n", countEmpresas)
	}

	// Contar estabelecimentos
	var countEstab int64
	err = db.QueryRow("SELECT COUNT(*) FROM receita.estabelecimento").Scan(&countEstab)
	if err != nil {
		log.Printf("⚠️  Erro ao contar estabelecimentos: %v", err)
	} else {
		fmt.Printf("✅ Estabelecimentos: %d\n", countEstab)
	}

	// Contar sócios
	var countSocios int64
	err = db.QueryRow("SELECT COUNT(*) FROM receita.socios").Scan(&countSocios)
	if err != nil {
		log.Printf("⚠️  Erro ao contar sócios: %v", err)
	} else {
		fmt.Printf("✅ Sócios: %d\n", countSocios)
	}

	// Testar query com TablePrefix
	fmt.Println("")
	fmt.Println("📋 Testando TablePrefix...")
	
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", database.TablePrefix("empresas"))
	fmt.Printf("Query: %s\n", query)
	
	var count int64
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("❌ Erro: %v", err)
	} else {
		fmt.Printf("✅ Resultado: %d\n", count)
	}

	// Testar AdaptQuery
	fmt.Println("")
	fmt.Println("🔄 Testando AdaptQuery...")
	
	sqliteQuery := "SELECT * FROM empresas WHERE cnpj_basico = ? AND razao_social LIKE ?"
	pgQuery := database.AdaptQuery(sqliteQuery)
	fmt.Printf("SQLite: %s\n", sqliteQuery)
	fmt.Printf("PostgreSQL: %s\n", pgQuery)

	fmt.Println("")
	fmt.Println("✅ Todos os testes concluídos!")
}
