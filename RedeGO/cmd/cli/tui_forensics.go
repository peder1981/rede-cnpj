package main

import (
	"fmt"
	"strings"

	"github.com/peder1981/rede-cnpj/RedeGO/internal/crossdata"
	"github.com/peder1981/rede-cnpj/RedeGO/internal/forensics"
)

// viewForensicsInvestigate exibe perfil completo de investigação
func (m model) viewForensicsInvestigate(cpf string) string {
	inv := forensics.NewInvestigator("bases/cnpj.db", "bases/rede.db")
	profile, err := inv.InvestigatePerson(cpf)
	
	if err != nil {
		return fmt.Sprintf("\n❌ ERRO: %v\n", err)
	}

	s := "\n"
	s += "╔══════════════════════════════════════════════════════════════════════╗\n"
	s += "║         🔍 INVESTIGAÇÃO FORENSE - PERFIL COMPLETO                   ║\n"
	s += "╚══════════════════════════════════════════════════════════════════════╝\n\n"

	// Dados pessoais
	s += "┌─ IDENTIFICAÇÃO ────────────────────────────────────────────────────┐\n"
	s += fmt.Sprintf("│ CPF:  %s                                                     │\n", profile.CPF)
	s += fmt.Sprintf("│ Nome: %-60s │\n", truncate(profile.Nome, 60))
	s += "└────────────────────────────────────────────────────────────────────┘\n\n"

	// Score de risco
	scoreBar := getScoreBar(profile.Score)
	scoreLevel := getScoreLevel(profile.Score)
	s += "┌─ SCORE DE RISCO ───────────────────────────────────────────────────┐\n"
	s += fmt.Sprintf("│ %s %3d/100 - %s%s │\n", scoreBar, profile.Score, scoreLevel, strings.Repeat(" ", 35-len(scoreLevel)))
	s += "└────────────────────────────────────────────────────────────────────┘\n\n"

	// Resumo
	s += "┌─ RESUMO EMPRESARIAL ───────────────────────────────────────────────┐\n"
	s += fmt.Sprintf("│ Total de Empresas:      %6d                                     │\n", profile.TotalEmpresas)
	s += fmt.Sprintf("│ • Ativas:               %6d                                     │\n", profile.EmpresasAtivas)
	s += fmt.Sprintf("│ • Baixadas:             %6d                                     │\n", profile.EmpresasBaixadas)
	s += fmt.Sprintf("│ • Suspensas:            %6d                                     │\n", profile.EmpresasSuspensas)
	s += fmt.Sprintf("│ Capital Social Total:   R$ %.2f                          │\n", profile.CapitalSocialTotal)
	s += "└────────────────────────────────────────────────────────────────────┘\n\n"

	// Indicadores
	s += "┌─ INDICADORES DE RISCO ─────────────────────────────────────────────┐\n"
	s += fmt.Sprintf("│ Endereços Diferentes:   %6d                                     │\n", profile.EnderecosDiferentes)
	s += fmt.Sprintf("│ Telefones Diferentes:   %6d                                     │\n", profile.TelefonesDiferentes)
	s += fmt.Sprintf("│ Emails Diferentes:      %6d                                     │\n", profile.EmailsDiferentes)
	s += fmt.Sprintf("│ Rede Bancária:          %6d empresas conectadas                 │\n", profile.RedeBancaria)
	s += "└────────────────────────────────────────────────────────────────────┘\n\n"

	// Flags
	if len(profile.Flags) > 0 {
		s += "┌─ ALERTAS ──────────────────────────────────────────────────────────┐\n"
		for _, flag := range profile.Flags {
			s += fmt.Sprintf("│ ⚠️  %-66s │\n", truncate(flag, 66))
		}
		s += "└────────────────────────────────────────────────────────────────────┘\n\n"
	}

	// Empresas (primeiras 10)
	if len(profile.Empresas) > 0 {
		s += "┌─ EMPRESAS (Top 10) ────────────────────────────────────────────────┐\n"
		for i, emp := range profile.Empresas {
			if i >= 10 {
				break
			}
			cnpj := emp["cnpj"].(string)
			razao := truncate(emp["razao_social"].(string), 40)
			situacao := emp["situacao"].(string)
			sitIcon := "✅"
			if situacao == "08" {
				sitIcon = "❌"
			} else if situacao != "02" {
				sitIcon = "⚠️"
			}
			s += fmt.Sprintf("│ %s %s - %-40s │\n", sitIcon, cnpj, razao)
		}
		if len(profile.Empresas) > 10 {
			s += fmt.Sprintf("│ ... e mais %d empresas                                            │\n", len(profile.Empresas)-10)
		}
		s += "└────────────────────────────────────────────────────────────────────┘\n"
	}

	s += "\n[Q] Voltar | [D] Ver Detalhes de Empresa | [T] Timeline\n"

	return s
}

// viewCPFDetails exibe todos os dados de um CPF
func (m model) viewCPFDetails(cpf string) string {
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	empresas, err := engine.EmpresasPorCPF(cpf)
	
	if err != nil {
		return fmt.Sprintf("\n❌ ERRO: %v\n", err)
	}

	s := "\n"
	s += "╔══════════════════════════════════════════════════════════════════════╗\n"
	s += "║         👤 DADOS COMPLETOS - CPF (SEM CENSURA)                      ║\n"
	s += "╚══════════════════════════════════════════════════════════════════════╝\n\n"

	if len(empresas) == 0 {
		s += "Nenhuma empresa encontrada para este CPF.\n"
		return s
	}

	// Dados da primeira empresa para pegar nome
	primeiro := empresas[0]
	
	s += "┌─ IDENTIFICAÇÃO ────────────────────────────────────────────────────┐\n"
	s += fmt.Sprintf("│ CPF:  %s (SEM MÁSCARA)                                      │\n", cpf)
	if nome, ok := primeiro["nome_socio"].(string); ok {
		s += fmt.Sprintf("│ Nome: %-60s │\n", truncate(nome, 60))
	}
	s += "└────────────────────────────────────────────────────────────────────┘\n\n"

	s += fmt.Sprintf("┌─ EMPRESAS (%d) ────────────────────────────────────────────────────┐\n", len(empresas))
	
	for i, emp := range empresas {
		if i >= 20 { // Limita a 20 para não poluir
			s += fmt.Sprintf("│ ... e mais %d empresas                                            │\n", len(empresas)-20)
			break
		}

		cnpj := getString(emp, "cnpj")
		razao := truncate(getString(emp, "razao_social"), 35)
		qualif := truncate(getString(emp, "qualificacao_socio"), 25)
		situacao := getString(emp, "situacao_cadastral")
		
		sitIcon := "✅"
		if situacao == "08" {
			sitIcon = "❌"
		} else if situacao != "02" {
			sitIcon = "⚠️"
		}

		s += fmt.Sprintf("│ %s %s                                                    │\n", sitIcon, cnpj)
		s += fmt.Sprintf("│    %-35s                                  │\n", razao)
		s += fmt.Sprintf("│    Cargo: %-25s                                │\n", qualif)
		
		// Dados de contato SEM CENSURA
		if email := getString(emp, "correio_eletronico"); email != "" {
			s += fmt.Sprintf("│    📧 %-60s │\n", truncate(email, 60))
		}
		if tel := getString(emp, "telefone1"); tel != "" {
			ddd := getString(emp, "ddd1")
			s += fmt.Sprintf("│    📞 (%s) %s                                                │\n", ddd, tel)
		}
		
		// Endereço SEM CENSURA
		logr := getString(emp, "logradouro")
		num := getString(emp, "numero")
		bairro := getString(emp, "bairro")
		cep := getString(emp, "cep")
		uf := getString(emp, "uf")
		
		if logr != "" {
			endereco := fmt.Sprintf("%s, %s - %s", logr, num, bairro)
			s += fmt.Sprintf("│    🏠 %-60s │\n", truncate(endereco, 60))
			s += fmt.Sprintf("│       %s - %s                                                  │\n", cep, uf)
		}
		
		s += "│                                                                      │\n"
	}
	
	s += "└────────────────────────────────────────────────────────────────────┘\n"
	s += "\n[Q] Voltar | [I] Investigar (Forense) | [T] Timeline\n"

	return s
}

// viewCNPJDetails exibe todos os dados de um CNPJ
func (m model) viewCNPJDetails(cnpj string) string {
	engine := crossdata.NewCrossDataEngine("bases/cnpj.db", "bases/rede.db")
	empresa, err := engine.DadosCompletosEmpresa(cnpj)
	
	if err != nil {
		return fmt.Sprintf("\n❌ ERRO: %v\n", err)
	}

	s := "\n"
	s += "╔══════════════════════════════════════════════════════════════════════╗\n"
	s += "║         🏢 DADOS COMPLETOS - CNPJ (SEM CENSURA)                     ║\n"
	s += "╚══════════════════════════════════════════════════════════════════════╝\n\n"

	// Identificação
	s += "┌─ IDENTIFICAÇÃO ────────────────────────────────────────────────────┐\n"
	s += fmt.Sprintf("│ CNPJ: %s                                                 │\n", empresa.CNPJ)
	s += fmt.Sprintf("│ Razão Social: %-55s │\n", truncate(empresa.RazaoSocial, 55))
	if empresa.NomeFantasia != "" {
		s += fmt.Sprintf("│ Nome Fantasia: %-54s │\n", truncate(empresa.NomeFantasia, 54))
	}
	s += fmt.Sprintf("│ Matriz/Filial: %-55s │\n", empresa.MatrizFilial)
	s += "└────────────────────────────────────────────────────────────────────┘\n\n"

	// Situação
	s += "┌─ SITUAÇÃO CADASTRAL ───────────────────────────────────────────────┐\n"
	s += fmt.Sprintf("│ Situação: %-59s │\n", empresa.SituacaoCadastral)
	s += fmt.Sprintf("│ Data: %s                                                     │\n", empresa.DataSituacaoCadastral)
	if empresa.MotivoSituacaoCadastral != "" {
		s += fmt.Sprintf("│ Motivo: %-61s │\n", truncate(empresa.MotivoSituacaoCadastral, 61))
	}
	s += fmt.Sprintf("│ Abertura: %s                                                 │\n", empresa.DataInicioAtividades)
	s += "└────────────────────────────────────────────────────────────────────┘\n\n"

	// Atividade
	s += "┌─ ATIVIDADE ECONÔMICA ──────────────────────────────────────────────┐\n"
	s += fmt.Sprintf("│ CNAE Principal: %-53s │\n", truncate(empresa.CNAEFiscal, 53))
	if empresa.CNAESecundaria != "" {
		s += fmt.Sprintf("│ CNAEs Secundárias: %-50s │\n", truncate(empresa.CNAESecundaria, 50))
	}
	s += fmt.Sprintf("│ Natureza Jurídica: %-51s │\n", truncate(empresa.NaturezaJuridica, 51))
	s += fmt.Sprintf("│ Porte: %-62s │\n", empresa.PorteEmpresa)
	s += fmt.Sprintf("│ Capital Social: R$ %.2f                                    │\n", empresa.CapitalSocial)
	if empresa.OpcaoMEI == "S" {
		s += "│ MEI: Sim                                                             │\n"
	}
	s += "└────────────────────────────────────────────────────────────────────┘\n\n"

	// ENDEREÇO COMPLETO - SEM CENSURA
	s += "┌─ ENDEREÇO (SEM CENSURA) ───────────────────────────────────────────┐\n"
	endereco := fmt.Sprintf("%s %s, %s", empresa.TipoLogradouro, empresa.Logradouro, empresa.Numero)
	s += fmt.Sprintf("│ %s%-68s │\n", "🏠 ", truncate(endereco, 68))
	if empresa.Complemento != "" {
		s += fmt.Sprintf("│    %-68s │\n", truncate(empresa.Complemento, 68))
	}
	s += fmt.Sprintf("│    %s - %s - %s                                           │\n", 
		empresa.Bairro, empresa.CEP, empresa.UF)
	s += fmt.Sprintf("│    Município: %-57s │\n", truncate(empresa.Municipio, 57))
	s += "└────────────────────────────────────────────────────────────────────┘\n\n"

	// CONTATOS - SEM CENSURA
	s += "┌─ CONTATOS (SEM CENSURA) ───────────────────────────────────────────┐\n"
	if empresa.Telefone1 != "" {
		s += fmt.Sprintf("│ 📞 Telefone 1: (%s) %s                                        │\n", 
			empresa.DDD1, empresa.Telefone1)
	}
	if empresa.Telefone2 != "" {
		s += fmt.Sprintf("│ 📞 Telefone 2: (%s) %s                                        │\n", 
			empresa.DDD2, empresa.Telefone2)
	}
	if empresa.Fax != "" {
		s += fmt.Sprintf("│ 📠 Fax: (%s) %s                                               │\n", 
			empresa.DDDFax, empresa.Fax)
	}
	if empresa.CorreioEletronico != "" {
		s += fmt.Sprintf("│ 📧 Email: %-61s │\n", truncate(empresa.CorreioEletronico, 61))
	}
	s += "└────────────────────────────────────────────────────────────────────┘\n"

	s += "\n[Q] Voltar | [S] Ver Sócios | [C] Cadeia de Controle\n"

	return s
}

// Funções auxiliares
func getScoreBar(score int) string {
	filled := score / 10
	empty := 10 - filled
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
}

func getScoreLevel(score int) string {
	if score >= 86 {
		return "🚨 CRÍTICO"
	} else if score >= 71 {
		return "🔴 ALTO RISCO"
	} else if score >= 51 {
		return "🔶 MÉDIO RISCO"
	} else if score >= 31 {
		return "⚠️  ATENÇÃO"
	}
	return "✅ BAIXO RISCO"
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
