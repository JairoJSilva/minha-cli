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
	Short:   "Escaneia e importa configurações existentes no terminal sem perder nada",
	Run: func(cmd *cobra.Command, args []string) {
		runScan()
	},
}

func runScan() {
	tui.PrintBanner()
	fmt.Println("\n🔍 Escaneando Configurações Existentes na Máquina...")

	res, err := config.ScanLocalEnvironment()
	if err != nil {
		tui.Error(fmt.Sprintf("Erro na varredura: %v", err))
		return
	}

	body := fmt.Sprintf(` AWS Profiles    : %d encontrados (%s)
 OCI Profiles    : %d encontrados (%s)
 GCP Configs     : %d encontrados (%s)
 K8s Contextos   : %d encontrados (%s)
 Azure Assinaturas: %d encontrados (%s)`,
		len(res.AWSProfiles), strings.Join(res.AWSProfiles, ", "),
		len(res.OCIProfiles), strings.Join(res.OCIProfiles, ", "),
		len(res.GCPConfigs), strings.Join(res.GCPConfigs, ", "),
		len(res.K8sContexts), strings.Join(res.K8sContexts, ", "),
		len(res.AzureSubs), strings.Join(res.AzureSubs, ", "),
	)

	tui.PrintCard("DIAGNÓSTICO DO TERMINAL / MÁQUINA", body)
	tui.Success("Varredura concluída com segurança!")
	fmt.Printf("  \033[32m✔ %d novo(s) cliente(s) importado(s) automaticamente.\033[0m\n", res.ImportedCount)
	fmt.Printf("  \033[36m✔ %d cliente(s) já existentes foram preservados intactos.\033[0m\n\n", res.ExistingCount)
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
