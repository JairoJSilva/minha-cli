package cmd

import (
	"fmt"
	"strings"

	"github.com/JairoJSilva/minha-cli/internal/config"
	"github.com/JairoJSilva/minha-cli/internal/env"
	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:     "scan",
	Aliases: []string{"leitura", "importar", "import"},
	Short:   "Escaneia e detalha todas as contas e perfis configurados no seu terminal",
	Run: func(cmd *cobra.Command, args []string) {
		runScan()
	},
}

func runScan() {
	tui.PrintBanner()
	fmt.Println("\n🔍 Escaneando Configurações Detalhadas no Terminal/Máquina...")

	report, err := config.ScanLocalEnvironmentDetailed()
	if err != nil {
		tui.Error(fmt.Sprintf("Erro na varredura: %v", err))
		return
	}

	// 1. Detalhes da AWS
	var awsLines strings.Builder
	if len(report.AWSDetails) > 0 {
		for _, a := range report.AWSDetails {
			awsLines.WriteString(fmt.Sprintf(" • Profile: \033[1m%s\033[0m | AccessKey: %s | Região: %s (%s)\n", a.Name, a.AccessKey, a.Region, a.Source))
		}
	} else {
		awsLines.WriteString(" • Nenhum profile encontrado em ~/.aws/credentials\n")
	}
	tui.PrintCard("☁️  AWS PROFILES DETECTADOS", strings.TrimRight(awsLines.String(), "\n"))

	// 2. Detalhes da Oracle OCI
	var ociLines strings.Builder
	if len(report.OCIDetails) > 0 {
		for _, o := range report.OCIDetails {
			ociLines.WriteString(fmt.Sprintf(" • Profile: \033[1m%s\033[0m | Tenancy: %s | Região: %s (~/.oci/config)\n", o.Name, o.Tenancy, o.Region))
		}
	} else {
		ociLines.WriteString(" • Nenhum profile encontrado em ~/.oci/config\n")
	}
	tui.PrintCard("🏛️  ORACLE CLOUD (OCI) PERFIS DETECTADOS", strings.TrimRight(ociLines.String(), "\n"))

	// 3. Detalhes do Kubernetes
	var k8sLines strings.Builder
	if len(report.K8sDetails) > 0 {
		for _, k := range report.K8sDetails {
			k8sLines.WriteString(fmt.Sprintf(" • Contexto: \033[1m%s\033[0m | Cluster: %s\n", k.Name, k.Cluster))
		}
	} else {
		k8sLines.WriteString(" • Nenhum cluster/contexto detectado no kubectl\n")
	}
	tui.PrintCard("☸️  KUBERNETES CONTEXTOS DETECTADOS", strings.TrimRight(k8sLines.String(), "\n"))

	// 4. Detalhes do Google Cloud e Azure
	var otherLines strings.Builder
	if len(report.GCPDetails) > 0 {
		for _, g := range report.GCPDetails {
			otherLines.WriteString(fmt.Sprintf(" • GCP Config: \033[1m%s\033[0m\n", g.Name))
		}
	}
	if len(report.AzureDetails) > 0 {
		for _, az := range report.AzureDetails {
			otherLines.WriteString(fmt.Sprintf(" • Azure Sub: \033[1m%s\033[0m (ID: %s)\n", az.Name, az.ID))
		}
	}
	if otherLines.Len() > 0 {
		tui.PrintCard("🌐 GOOGLE CLOUD & AZURE DETECTADOS", strings.TrimRight(otherLines.String(), "\n"))
	}

	tui.Success("Varredura concluída com proteção total de dados!")
	fmt.Printf("  \033[32m✔ %d novo(s) cliente(s) importado(s) automaticamente para a CLI.\033[0m\n", report.ImportedCount)
	fmt.Printf("  \033[36m✔ %d perfil(is) já existentes foram preservados intactos.\033[0m\n\n", report.ExistingCount)
}

var clearCmd = &cobra.Command{
	Use:     "clear",
	Aliases: []string{"reset", "c", "limpar"},
	Short:   "Limpa todas as variáveis de ambiente das nuvens",
	Run: func(cmd *cobra.Command, args []string) {
		runClear()
	},
}

func runClear() {
	env.WriteEnvToFile(env.GenerateClearScript())
	tui.Success("Todas as variáveis de ambiente das nuvens foram limpas com sucesso.")
}

var envCmd = &cobra.Command{
	Use:    "env [cliente]",
	Hidden: true,
	Short:  "Gera os comandos shell de export/unset para avaliação do wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "clear" {
			fmt.Print(env.GenerateClearScript())
			return
		}

		c, err := config.FindClient(args[0])
		if err == nil && c != nil {
			fmt.Print(env.GenerateExportScript(c))
		}
	},
}
