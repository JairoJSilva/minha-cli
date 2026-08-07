package cmd

import (
	"fmt"

	"github.com/JairoJSilva/minha-cli/internal/config"
	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:     "switch [cliente]",
	Aliases: []string{"s", "trocar"},
	Short:   "Alterna o contexto multi-cloud e Kubernetes para o cliente informado",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			applyClientDirect(args[0])
		} else {
			runSwitchInteractive()
		}
	},
}

func applyClientDirect(target string) {
	client, err := config.FindClient(target)
	if err != nil {
		tui.Error(fmt.Sprintf("Cliente '%s' não encontrado. Use 'mc list' para ver os disponíveis.", target))
		return
	}

	tui.Success(fmt.Sprintf("Contexto ativado: %s", client.Name))
	aws := "<nenhum>"
	if client.AWSProfile != nil {
		aws = *client.AWSProfile
	}
	oci := "<nenhum>"
	if client.OCIProfile != nil {
		oci = *client.OCIProfile
	}
	gcp := "<nenhum>"
	if client.GCPConfig != nil {
		gcp = *client.GCPConfig
	}
	k8s := "<nenhum>"
	if client.K8sContext != nil {
		k8s = *client.K8sContext
	}

	fmt.Printf("  \033[2mAWS: %s | OCI: %s | GCP: %s | K8s: %s\033[0m\n", aws, oci, gcp, k8s)
}

func runSwitchInteractive() {
	clients, err := config.LoadClients()
	if err != nil || len(clients) == 0 {
		tui.Warn("Nenhum cliente cadastrado. Use 'mc add' para cadastrar o primeiro.")
		return
	}

	var options []huh.Option[string]
	for _, c := range clients {
		options = append(options, huh.NewOption(c.Name, c.ID))
	}
	options = append(options, huh.NewOption("🧹 Limpar Contexto (Reset)", "clear"))

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Selecione o Cliente / Contexto").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return
	}

	if selected == "clear" {
		runClear()
		return
	}

	applyClientDirect(selected)
}
