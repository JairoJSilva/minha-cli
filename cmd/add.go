package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/JairoJSilva/minha-cli/internal/config"
	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"novo", "criar", "new"},
	Short:   "Cadastra uma nova conta / cliente interativamente",
	Run: func(cmd *cobra.Command, args []string) {
		runAddInteractive()
	},
}

func runAddInteractive() {
	tui.PrintBanner()
	fmt.Println("\n➕ Cadastrar Nova Conta / Cliente")

	var name, id, aws, oci, gcp, azure, k8s string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Nome de Exibição").
				Description("Ex: Hospital Albert Einstein, Santander").
				Placeholder("Nome do Cliente").
				Value(&name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("o nome não pode ser vazio")
					}
					return nil
				}),

			huh.NewInput().
				Title("ID Curto / Slug").
				Description("Identificador único para usar no 'mc switch <id>'").
				Placeholder("ex: einstein").
				Value(&id),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("AWS Profile (Opcional)").
				Placeholder("ex: einstein-prod").
				Value(&aws),

			huh.NewInput().
				Title("Oracle OCI Profile (Opcional)").
				Placeholder("ex: einstein-oci").
				Value(&oci),

			huh.NewInput().
				Title("Google Cloud Config (Opcional)").
				Placeholder("ex: einstein-gcp").
				Value(&gcp),

			huh.NewInput().
				Title("Azure Subscription ID (Opcional)").
				Placeholder("ex: xxxxxxxx-xxxx-xxxx").
				Value(&azure),

			huh.NewInput().
				Title("Contexto Kubernetes (Opcional)").
				Placeholder("ex: cluster-einstein").
				Value(&k8s),
		),
	)

	if err := form.Run(); err != nil {
		tui.Warn("Cadastro cancelado.")
		return
	}

	if id == "" {
		id = strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(name, ""))
	}

	client := config.Client{
		ID:         id,
		Name:       name,
		AWSProfile: config.StringPtr(aws),
		OCIProfile: config.StringPtr(oci),
		GCPConfig:  config.StringPtr(gcp),
		AzureSub:   config.StringPtr(azure),
		K8sContext: config.StringPtr(k8s),
	}

	if err := config.AddClient(client); err != nil {
		tui.Error(fmt.Sprintf("Falha ao salvar cliente: %v", err))
		return
	}

	tui.Success(fmt.Sprintf("Cliente '%s' (%s) adicionado com sucesso!", name, id))
}
