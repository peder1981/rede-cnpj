package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config representa a configuração da aplicação
type Config struct {
	// Base de dados
	BaseReceita             string
	BaseRede                string
	BaseRedeSearch          string
	BaseEnderecoNormalizado string
	BaseLinks               string
	BaseLocal               string
	PastaArquivos           string
	ReferenciaBD            string

	// Servidor
	PortaFlask int

	// Limites
	LimiterPadrao   string
	LimiterDados    string
	LimiterGoogle   string
	LimiterArquivos string

	// Flags
	BuscaGoogle         bool
	BuscaChaves         bool
	ArquivosDownload    bool
	LigacaoSocioFilial  bool
	ExibeMensagemInicial bool
	ExibeMenuInserir    bool

	// Limites de processamento
	LimiteRegistrosCamada int
	TempoMaximoConsulta   float64
	GeocodeMax            int

	// API
	APICnpj     bool
	APICaminhos bool
	APIKeys     []string

	// Parâmetros de linha de comando
	CPFCNPJInicial       string
	CamadaInicial        int
	IDArquivoServidor    string
	ArquivoEntrada       string
	EncodingArquivo      string
	BTextoEmbaixoIcone   bool
	ExcelSheetName       interface{}
	Separador            string
	TipoLista            string
}

var AppConfig *Config

// LoadConfig carrega a configuração do arquivo INI
func LoadConfig() (*Config, error) {
	// Parse command line flags
	confFile := flag.String("conf_file", "rede.ini", "arquivo de configuração")
	cpfcnpjInicial := flag.String("inicial", "", "1 ou mais cnpj separados por ponto e vírgula")
	camadaInicial := flag.Int("camada", 1, "camada inicial")
	idArquivoServidor := flag.String("json", "", "nome json no servidor")
	arquivoEntrada := flag.String("lista", "", "inserir itens de arquivo em gráfico")
	encodingArquivo := flag.String("encoding", "utf8", "codificação do arquivo")
	pastaArquivos := flag.String("pasta", "arquivos", "pasta de arquivos do usuário")
	portaFlask := flag.Int("porta_flask", 5000, "porta da aplicação")
	textoEmbaixo := flag.Bool("texto-embaixo", true, "texto embaixo do ícone")
	mensagem := flag.Bool("mensagem", true, "exibe mensagem inicial")
	menuInserir := flag.Bool("menuinserir", true, "exibe menu para inserir no inicio")
	download := flag.Bool("download", false, "permitir download da pasta arquivos")
	sheetName := flag.String("sheet-name", "0", "nome da aba do excel")
	separador := flag.String("separador", "\t", "separador arquivo csv")
	tipoLista := flag.String("tipo_lista", "", "define tipo de entrada")

	flag.Parse()

	// Verifica se o arquivo de configuração existe
	if _, err := os.Stat(*confFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("arquivo de configuração %s não encontrado", *confFile)
	}

	// Configura Viper para ler arquivo INI
	viper.SetConfigFile(*confFile)
	viper.SetConfigType("ini")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
	}

	cfg := &Config{
		// Base de dados
		BaseReceita:             viper.GetString("BASE.base_receita"),
		BaseRede:                viper.GetString("BASE.base_rede"),
		BaseRedeSearch:          viper.GetString("BASE.base_rede_search"),
		BaseEnderecoNormalizado: viper.GetString("BASE.base_endereco_normalizado"),
		BaseLinks:               viper.GetString("BASE.base_links"),
		BaseLocal:               viper.GetString("BASE.base_local"),
		ReferenciaBD:            viper.GetString("BASE.referencia_bd"),

		// Servidor
		PortaFlask: *portaFlask,

		// Limites
		LimiterPadrao:   viper.GetString("ETC.limiter_padrao"),
		LimiterDados:    viper.GetString("ETC.limiter_dados"),
		LimiterGoogle:   viper.GetString("ETC.limiter_google"),
		LimiterArquivos: viper.GetString("ETC.limiter_arquivos"),

		// Flags
		BuscaGoogle:          viper.GetBool("ETC.busca_google"),
		BuscaChaves:          viper.GetBool("ETC.busca_chaves"),
		LigacaoSocioFilial:   viper.GetBool("ETC.ligacao_socio_filial"),
		ExibeMensagemInicial: *mensagem,
		ExibeMenuInserir:     *menuInserir,
		ArquivosDownload:     *download,

		// Limites de processamento
		LimiteRegistrosCamada: viper.GetInt("ETC.limite_registros_camada"),
		TempoMaximoConsulta:   viper.GetFloat64("ETC.tempo_maximo_consulta"),
		GeocodeMax:            viper.GetInt("ETC.geocode_max"),

		// API
		APICnpj:     viper.GetBool("API.api_cnpj"),
		APICaminhos: viper.GetBool("API.api_caminhos"),
		APIKeys:     viper.GetStringSlice("API.api_keys"),

		// Parâmetros de linha de comando
		CPFCNPJInicial:     *cpfcnpjInicial,
		CamadaInicial:      *camadaInicial,
		IDArquivoServidor:  *idArquivoServidor,
		ArquivoEntrada:     *arquivoEntrada,
		EncodingArquivo:    *encodingArquivo,
		PastaArquivos:      *pastaArquivos,
		BTextoEmbaixoIcone: *textoEmbaixo,
		ExcelSheetName:     *sheetName,
		Separador:          *separador,
		TipoLista:          *tipoLista,
	}

	// Valores padrão
	if cfg.LimiterPadrao == "" {
		cfg.LimiterPadrao = "20/minute"
	}
	if cfg.LimiterDados == "" {
		cfg.LimiterDados = cfg.LimiterPadrao
	}
	if cfg.LimiterGoogle == "" {
		cfg.LimiterGoogle = "4/minute"
	}
	if cfg.LimiterArquivos == "" {
		cfg.LimiterArquivos = "2/minute"
	}
	if cfg.LimiteRegistrosCamada == 0 {
		cfg.LimiteRegistrosCamada = 1000
	}
	if cfg.TempoMaximoConsulta == 0 {
		cfg.TempoMaximoConsulta = 10.0
	}
	if cfg.GeocodeMax == 0 {
		cfg.GeocodeMax = 15
	}
	if cfg.PastaArquivos == "" {
		cfg.PastaArquivos = "arquivos"
	}

	// Verifica se os arquivos de banco de dados existem
	databases := map[string]string{
		"base_receita":    cfg.BaseReceita,
		"base_rede":       cfg.BaseRede,
		"base_rede_search": cfg.BaseRedeSearch,
	}

	for name, path := range databases {
		if path != "" {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Printf("AVISO: Arquivo %s (%s) não encontrado\n", name, path)
			}
		}
	}

	AppConfig = cfg
	return cfg, nil
}

// GetConfig retorna a configuração global
func GetConfig() *Config {
	if AppConfig == nil {
		panic("configuração não foi carregada. Chame LoadConfig() primeiro")
	}
	return AppConfig
}
