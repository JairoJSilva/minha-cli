package cmd

import (
	"fmt"
	"strings"

	"github.com/JairoJSilva/minha-cli/internal/config"
	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:     "show [cliente]",
	Aliases: []string{"detalhes", "info", "view"},
	Short:   "Exibe os detalhes completos de uma conta/cliente cadastrado",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			showClientDetails(args[0])
		} else {
			runShowInteractive()
		}
	},
}

func showClientDetails(target string) {
	client, err := config.FindClient(target)
	if err != nil {
		tui.Error(fmt.Sprintf("Cliente '%s' não encontrado. Use 'mc list' para ver os cadastrados.", target))
		return
	}

	aws := config.SafeString(client.AWSProfile)
	awsStatus := "Não configurado"
	if aws != "" {
		if config.AWSProfileLocalExists(aws) {
			awsStatus = "Configurado em ~/.aws/credentials"
		} else {
			awsStatus = "⚠️ Profile ausente no ~/.aws/credentials"
		}
	} else {
		aws = "-"
	}

	oci := config.SafeString(client.OCIProfile)
	ociStatus := "Não configurado"
	if oci != "" {
		ociStatus = "Configurado em ~/.oci/config"
	} else {
		oci = "-"
	}

	gcp := config.SafeString(client.GCPConfig)
	if gcp == "" {
		gcp = "-"
	}

	azure := config.SafeString(client.AzureSub)
	if azure == "" {
		azure = "-"
	}

	k8s := config.SafeString(client.K8sContext)
	if k8s == "" {
		k8s = "-"
	}

	body := fmt.Sprintf(` 🏷️  Identificador : %s
 📝 Nome Completo : %s

 ☁️  AWS Cloud:
    └─ Profile    : %s
    └─ Status     : %s

 🏛️  Oracle Cloud (OCI):
    └─ Profile    : %s
    └─ Status     : %s

 🌐 Google Cloud  : %s
 🔷 Azure Sub     : %s
 ☸️  Kubernetes    : %s`,
		client.ID, client.Name, aws, awsStatus, oci, ociStatus, gcp, azure, k8s)

	tui.PrintBanner()
	tui.PrintCard(fmt.Sprintf("DETALHES DA CONTA — %s", strings.ToUpper(client.ID)), body)
}

func runShowInteractive() {
	clients, err := config.LoadClients()
	if err != nil || len(clients) == 0 {
		tui.Warn("Nenhum cliente cadastrado. Use 'mc add' ou 'mc scan'.")
		return
	}

	var options []huh.Option[string]
	for _, c := range clients {
		options = append(options, huh.NewOption(fmt.Sprintf("%s (%s)", c.Name, c.ID), c.ID))
	}

	var selectedID string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Selecione a conta para ver os DETALHES").
				Options(options...).
				Value(&selectedID),
		),
	)

	if err := form.Run(); err != nil || selectedID == "" {
		return
	}

	showClientDetails(selectedID)
}
