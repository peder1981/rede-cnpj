package importer

import (
	"github.com/peder1981/rede-cnpj/RedeGO/internal/config"
)

// Importer gerencia a importação de dados da Receita Federal
type Importer struct {
	cfg *config.Config
	
	// Diretórios
	zipDir    string
	csvDir    string
	dbDir     string
	
	// URLs
	baseURL   string
}

// NewImporter cria um novo importador
func NewImporter(cfg *config.Config) *Importer {
	return &Importer{
		cfg:     cfg,
		zipDir:  "dados-publicos-zip",
		csvDir:  "dados-publicos",
		dbDir:   "bases",
		baseURL: "https://arquivos.receitafederal.gov.br/dados/cnpj/dados_abertos_cnpj/",
	}
}

// DownloadFiles baixa os arquivos ZIP da Receita Federal
func (i *Importer) DownloadFiles() error {
	downloader := NewDownloader(i.baseURL, i.zipDir)
	return downloader.Download()
}

// ProcessFiles processa os arquivos ZIP e cria o banco cnpj.db
func (i *Importer) ProcessFiles() error {
	processor := NewProcessorWithConfig(i.zipDir, i.csvDir, i.dbDir, i.cfg)
	return processor.Process()
}

// CreateLinkTables cria as tabelas de ligação (rede.db)
func (i *Importer) CreateLinkTables() error {
	linker := NewLinker(i.dbDir)
	return linker.CreateLinks()
}

// CreateSearchIndexes cria os índices de busca (rede_search.db)
func (i *Importer) CreateSearchIndexes() error {
	indexer := NewIndexer(i.dbDir)
	return indexer.CreateIndexes()
}
