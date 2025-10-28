-- ============================================================================
-- RedeCNPJ - PostgreSQL Tables Creation
-- ============================================================================
-- Cria as tabelas principais com particionamento
-- ============================================================================

-- ============================================================================
-- SCHEMA: receita (Dados da Receita Federal)
-- ============================================================================

-- Tabela: empresas (dados cadastrais básicos)
CREATE TABLE receita.empresas (
    cnpj_basico VARCHAR(8) PRIMARY KEY,
    razao_social TEXT NOT NULL,
    natureza_juridica VARCHAR(4),
    qualificacao_responsavel VARCHAR(2),
    capital_social NUMERIC(15,2) DEFAULT 0,
    porte_empresa VARCHAR(2),
    ente_federativo_responsavel TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE receita.empresas IS 'Dados cadastrais básicos das empresas';
COMMENT ON COLUMN receita.empresas.cnpj_basico IS '8 primeiros dígitos do CNPJ';
COMMENT ON COLUMN receita.empresas.razao_social IS 'Razão social da empresa';
COMMENT ON COLUMN receita.empresas.capital_social IS 'Capital social em reais';

-- Tabela: estabelecimento (matriz e filiais) - PARTICIONADA POR UF
CREATE TABLE receita.estabelecimento (
    cnpj VARCHAR(14) NOT NULL,
    cnpj_basico VARCHAR(8) NOT NULL,
    cnpj_ordem VARCHAR(4) NOT NULL,
    cnpj_dv VARCHAR(2) NOT NULL,
    matriz_filial VARCHAR(1),
    nome_fantasia TEXT,
    situacao_cadastral VARCHAR(2),
    data_situacao_cadastral DATE,
    motivo_situacao_cadastral VARCHAR(2),
    nome_cidade_exterior TEXT,
    pais VARCHAR(3),
    data_inicio_atividades DATE,
    cnae_fiscal VARCHAR(7),
    cnae_fiscal_secundaria TEXT,
    tipo_logradouro TEXT,
    logradouro TEXT,
    numero TEXT,
    complemento TEXT,
    bairro TEXT,
    cep VARCHAR(8),
    uf VARCHAR(2) NOT NULL,
    municipio VARCHAR(4),
    ddd1 VARCHAR(4),
    telefone1 VARCHAR(8),
    ddd2 VARCHAR(4),
    telefone2 VARCHAR(8),
    ddd_fax VARCHAR(4),
    fax VARCHAR(8),
    correio_eletronico TEXT,
    situacao_especial TEXT,
    data_situacao_especial DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (cnpj, uf),
    FOREIGN KEY (cnpj_basico) REFERENCES receita.empresas(cnpj_basico)
) PARTITION BY LIST (uf);

COMMENT ON TABLE receita.estabelecimento IS 'Estabelecimentos (matriz e filiais) particionados por UF';

-- Criar partições por UF (27 estados + DF)
CREATE TABLE receita.estabelecimento_ac PARTITION OF receita.estabelecimento FOR VALUES IN ('AC');
CREATE TABLE receita.estabelecimento_al PARTITION OF receita.estabelecimento FOR VALUES IN ('AL');
CREATE TABLE receita.estabelecimento_ap PARTITION OF receita.estabelecimento FOR VALUES IN ('AP');
CREATE TABLE receita.estabelecimento_am PARTITION OF receita.estabelecimento FOR VALUES IN ('AM');
CREATE TABLE receita.estabelecimento_ba PARTITION OF receita.estabelecimento FOR VALUES IN ('BA');
CREATE TABLE receita.estabelecimento_ce PARTITION OF receita.estabelecimento FOR VALUES IN ('CE');
CREATE TABLE receita.estabelecimento_df PARTITION OF receita.estabelecimento FOR VALUES IN ('DF');
CREATE TABLE receita.estabelecimento_es PARTITION OF receita.estabelecimento FOR VALUES IN ('ES');
CREATE TABLE receita.estabelecimento_go PARTITION OF receita.estabelecimento FOR VALUES IN ('GO');
CREATE TABLE receita.estabelecimento_ma PARTITION OF receita.estabelecimento FOR VALUES IN ('MA');
CREATE TABLE receita.estabelecimento_mt PARTITION OF receita.estabelecimento FOR VALUES IN ('MT');
CREATE TABLE receita.estabelecimento_ms PARTITION OF receita.estabelecimento FOR VALUES IN ('MS');
CREATE TABLE receita.estabelecimento_mg PARTITION OF receita.estabelecimento FOR VALUES IN ('MG');
CREATE TABLE receita.estabelecimento_pa PARTITION OF receita.estabelecimento FOR VALUES IN ('PA');
CREATE TABLE receita.estabelecimento_pb PARTITION OF receita.estabelecimento FOR VALUES IN ('PB');
CREATE TABLE receita.estabelecimento_pr PARTITION OF receita.estabelecimento FOR VALUES IN ('PR');
CREATE TABLE receita.estabelecimento_pe PARTITION OF receita.estabelecimento FOR VALUES IN ('PE');
CREATE TABLE receita.estabelecimento_pi PARTITION OF receita.estabelecimento FOR VALUES IN ('PI');
CREATE TABLE receita.estabelecimento_rj PARTITION OF receita.estabelecimento FOR VALUES IN ('RJ');
CREATE TABLE receita.estabelecimento_rn PARTITION OF receita.estabelecimento FOR VALUES IN ('RN');
CREATE TABLE receita.estabelecimento_rs PARTITION OF receita.estabelecimento FOR VALUES IN ('RS');
CREATE TABLE receita.estabelecimento_ro PARTITION OF receita.estabelecimento FOR VALUES IN ('RO');
CREATE TABLE receita.estabelecimento_rr PARTITION OF receita.estabelecimento FOR VALUES IN ('RR');
CREATE TABLE receita.estabelecimento_sc PARTITION OF receita.estabelecimento FOR VALUES IN ('SC');
CREATE TABLE receita.estabelecimento_sp PARTITION OF receita.estabelecimento FOR VALUES IN ('SP');
CREATE TABLE receita.estabelecimento_se PARTITION OF receita.estabelecimento FOR VALUES IN ('SE');
CREATE TABLE receita.estabelecimento_to PARTITION OF receita.estabelecimento FOR VALUES IN ('TO');
CREATE TABLE receita.estabelecimento_ex PARTITION OF receita.estabelecimento FOR VALUES IN ('EX'); -- Exterior

-- Tabela: socios - PARTICIONADA POR TIPO
CREATE TABLE receita.socios (
    id BIGSERIAL,
    cnpj VARCHAR(14) NOT NULL,
    cnpj_basico VARCHAR(8) NOT NULL,
    identificador_de_socio VARCHAR(1) NOT NULL,
    nome_socio TEXT NOT NULL,
    cnpj_cpf_socio VARCHAR(14) NOT NULL,
    qualificacao_socio VARCHAR(2),
    data_entrada_sociedade DATE,
    pais VARCHAR(3),
    representante_legal VARCHAR(11),
    nome_representante TEXT,
    qualificacao_representante_legal VARCHAR(2),
    faixa_etaria VARCHAR(1),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id, identificador_de_socio)
    -- Nota: FOREIGN KEY removida devido à complexidade com particionamento
    -- A integridade referencial deve ser garantida pela aplicação
) PARTITION BY LIST (identificador_de_socio);

COMMENT ON TABLE receita.socios IS 'Sócios e administradores particionados por tipo';
COMMENT ON COLUMN receita.socios.identificador_de_socio IS '1=PF, 2=PJ, 3=Estrangeiro';
COMMENT ON COLUMN receita.socios.cnpj_cpf_socio IS 'CPF ou CNPJ do sócio (SEM MÁSCARA)';

-- Criar partições por tipo de sócio
CREATE TABLE receita.socios_pf PARTITION OF receita.socios FOR VALUES IN ('1');
CREATE TABLE receita.socios_pj PARTITION OF receita.socios FOR VALUES IN ('2');
CREATE TABLE receita.socios_pe PARTITION OF receita.socios FOR VALUES IN ('3');

COMMENT ON TABLE receita.socios_pf IS 'Sócios Pessoa Física';
COMMENT ON TABLE receita.socios_pj IS 'Sócios Pessoa Jurídica';
COMMENT ON TABLE receita.socios_pe IS 'Sócios Estrangeiros';

-- Tabela: simples (Simples Nacional e MEI)
CREATE TABLE receita.simples (
    cnpj_basico VARCHAR(8) PRIMARY KEY,
    opcao_simples VARCHAR(1),
    data_opcao_simples DATE,
    data_exclusao_simples DATE,
    opcao_mei VARCHAR(1),
    data_opcao_mei DATE,
    data_exclusao_mei DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (cnpj_basico) REFERENCES receita.empresas(cnpj_basico)
);

COMMENT ON TABLE receita.simples IS 'Opção pelo Simples Nacional e MEI';

-- ============================================================================
-- Tabelas de Códigos (Lookup Tables)
-- ============================================================================

CREATE TABLE receita.cnae (
    codigo VARCHAR(7) PRIMARY KEY,
    descricao TEXT NOT NULL
);
COMMENT ON TABLE receita.cnae IS 'Classificação Nacional de Atividades Econômicas';

CREATE TABLE receita.motivo (
    codigo VARCHAR(2) PRIMARY KEY,
    descricao TEXT NOT NULL
);
COMMENT ON TABLE receita.motivo IS 'Motivos de situação cadastral';

CREATE TABLE receita.municipio (
    codigo VARCHAR(4) PRIMARY KEY,
    descricao TEXT NOT NULL
);
COMMENT ON TABLE receita.municipio IS 'Municípios brasileiros';

CREATE TABLE receita.natureza_juridica (
    codigo VARCHAR(4) PRIMARY KEY,
    descricao TEXT NOT NULL
);
COMMENT ON TABLE receita.natureza_juridica IS 'Natureza jurídica das empresas';

CREATE TABLE receita.pais (
    codigo VARCHAR(3) PRIMARY KEY,
    descricao TEXT NOT NULL
);
COMMENT ON TABLE receita.pais IS 'Países';

CREATE TABLE receita.qualificacao_socio (
    codigo VARCHAR(2) PRIMARY KEY,
    descricao TEXT NOT NULL
);
COMMENT ON TABLE receita.qualificacao_socio IS 'Qualificação de sócios';

-- ============================================================================
-- SCHEMA: rede (Relacionamentos)
-- ============================================================================

CREATE TABLE rede.ligacao (
    id BIGSERIAL PRIMARY KEY,
    id1 TEXT NOT NULL,
    id2 TEXT NOT NULL,
    descricao TEXT NOT NULL,
    cnpj VARCHAR(14),
    peso INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE rede.ligacao IS 'Relacionamentos entre entidades';
COMMENT ON COLUMN rede.ligacao.id1 IS 'ID origem (PJ_cnpj, PF_cpf-nome, PE_nome)';
COMMENT ON COLUMN rede.ligacao.id2 IS 'ID destino';
COMMENT ON COLUMN rede.ligacao.descricao IS 'Tipo de relacionamento';

-- ============================================================================
-- SCHEMA: forensics (Análises Forenses)
-- ============================================================================

CREATE TABLE forensics.perfil_risco (
    cpf VARCHAR(11) PRIMARY KEY,
    nome TEXT NOT NULL,
    score INTEGER NOT NULL CHECK (score >= 0 AND score <= 100),
    total_empresas INTEGER DEFAULT 0,
    empresas_ativas INTEGER DEFAULT 0,
    empresas_baixadas INTEGER DEFAULT 0,
    empresas_suspensas INTEGER DEFAULT 0,
    capital_social_total NUMERIC(15,2) DEFAULT 0,
    rede_bancaria INTEGER DEFAULT 0,
    flags JSONB,
    ultima_atualizacao TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE forensics.perfil_risco IS 'Perfis de risco calculados';

CREATE TABLE forensics.cluster_empresas (
    id SERIAL PRIMARY KEY,
    tipo_cluster VARCHAR(50) NOT NULL,
    criterio TEXT NOT NULL,
    valor_comum TEXT NOT NULL,
    total_empresas INTEGER NOT NULL,
    score INTEGER NOT NULL,
    flags JSONB,
    empresas JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE forensics.cluster_empresas IS 'Clusters de empresas suspeitas';

CREATE TABLE forensics.alertas (
    id BIGSERIAL PRIMARY KEY,
    tipo VARCHAR(50) NOT NULL,
    severidade VARCHAR(20) NOT NULL,
    entidade_id TEXT NOT NULL,
    entidade_tipo VARCHAR(10) NOT NULL,
    descricao TEXT NOT NULL,
    dados JSONB,
    status VARCHAR(20) DEFAULT 'NOVO',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE forensics.alertas IS 'Alertas de atividades suspeitas';

-- Criar triggers para updated_at (apenas em tabelas não particionadas)
CREATE TRIGGER update_empresas_updated_at BEFORE UPDATE ON receita.empresas
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Nota: Tabelas particionadas (estabelecimento, socios) não suportam triggers na tabela pai
-- Os triggers devem ser criados em cada partição individualmente se necessário

CREATE TRIGGER update_simples_updated_at BEFORE UPDATE ON receita.simples
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_alertas_updated_at BEFORE UPDATE ON forensics.alertas
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Mensagem de sucesso
DO $$
BEGIN
    RAISE NOTICE 'Tabelas criadas com sucesso!';
    RAISE NOTICE 'Partições criadas: 28 UFs + 3 tipos de sócios';
END $$;
