package crossdata

import "database/sql"

// CrossDataService fornece funcionalidades de cruzamento de dados
type CrossDataService struct {
	db *sql.DB
}

// NewCrossDataService cria um novo serviço de cruzamento
func NewCrossDataService(db *sql.DB) *CrossDataService {
	return &CrossDataService{db: db}
}
