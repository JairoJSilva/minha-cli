package cmd

import (
	"fmt"

	"github.com/JairoJSilva/minha-cli/internal/config"
	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     "edit [cliente]",
	Aliases: []string{"editar", "update"},
	Short:   "Edita configurações de um cliente existente",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			runEditClientByID(args[0])
		} else {
			runEditInteractive()
		}
	},
}

func runEditClientByID(selectedID string) {
	client, err := config.FindClient(selectedID)
	if err != nil {
		tui.Error(fmt.Sprintf("Cliente '%s' não encontrado.", selectedID))
		return
	}

	name := client.Name
	aws := config.SafeString(client.AWSProfile)
	oci := config.SafeString(client.OCIProfile)
	gcp := config.SafeString(client.GCPConfig)
	azure := config.SafeString(client.AzureSub)
	k8s := config.SafeString(client.K8sContext)

	var configAWSKeys bool
	var awsAccessKey, awsSecretKey string

	formEdit := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Nome de Exibição").
				Value(&name),
			huh.NewInput().
				Title("AWS Profile").
				Value(&aws),
			huh.NewInput().
				Title("Oracle OCI Profile").
				Value(&oci),
			huh.NewInput().
				Title("Google Cloud Config").
				Value(&gcp),
			huh.NewInput().
				Title("Azure Subscription ID").
				Value(&azure),
			huh.NewInput().
				Title("Contexto Kubernetes").
				Value(&k8s),
		),
	)

	if err := formEdit.Run(); err != nil {
		tui.Warn("Edição cancelada.")
		return
	}

	// Se o AWS profile preenchido não existir no ~/.aws/credentials, oferece gravar
	if aws != "" && !config.AWSProfileLocalExists(aws) {
		formAWS := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("O profile AWS '%s' ainda não existe no ~/.aws/credentials. Deseja cadastrar as chaves agora?", aws)).
					Value(&configAWSKeys),
			),
		)
		_ = formAWS.Run()

		if configAWSKeys {
			formKeys := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("AWS Access Key ID").
						Placeholder("AKIA...").
						Value(&awsAccessKey),
					huh.NewInput().
						Title("AWS Secret Access Key").
						EchoMode(huh.EchoModePassword).
						Value(&awsSecretKey),
				),
			)
			_ = formKeys.Run()
			if awsAccessKey != "" && awsSecretKey != "" {
				_ = saveAWSCredential(aws, awsAccessKey, awsSecretKey)
				tui.Success(fmt.Sprintf("Chaves salvas no ~/.aws/credentials para o perfil '%s'!", aws))
			}
		}
	}

	updated := config.Client{
		ID:         client.ID,
		Name:       name,
		AWSProfile: config.StringPtr(aws),
		OCIProfile: config.StringPtr(oci),
		GCPConfig:  config.StringPtr(gcp),
		AzureSub:   config.StringPtr(azure),
		K8sContext: config.StringPtr(k8s),
	}

	if err := config.UpdateClient(updated); err != nil {
		tui.Error(fmt.Sprintf("Falha ao atualizar: %v", err))
		return
	}

	tui.Success(fmt.Sprintf("Cliente '%s' atualizado com sucesso!", name))
}

func runEditInteractive() {
	clients, err := config.LoadClients()
	if err != nil || len(clients) == 0 {
		tui.Warn("Nenhum cliente cadastrado para editar.")
		return
	}

	var options []huh.Option[string]
	for _, c := range clients {
		options = append(options, huh.NewOption(c.Name, c.ID))
	}

	var selectedID string
	formSelect := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Selecione o cliente para EDITAR").
				Options(options...).
				Value(&selectedID),
		),
	)

	if err := formSelect.Run(); err != nil || selectedID == "" {
		return
	}

	runEditClientByID(selectedID)
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"apagar", "remover", "rm", "del"},
	Short:   "Remove uma conta/cliente cadastrado com segurança",
	Run: func(cmd *cobra.Command, args []string) {
		runDeleteInteractive()
	},
}

func runDeleteInteractive() {
	clients, err := config.LoadClients()
	if err != nil || len(clients) == 0 {
		tui.Warn("Nenhum cliente cadastrado.")
		return
	}

	var options []huh.Option[string]
	for _, c := range clients {
		options = append(options, huh.NewOption(c.Name, c.ID))
	}

	var selectedID string
	var confirmDelete bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Selecione o cliente para APAGAR").
				Options(options...).
				Value(&selectedID),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Tem certeza que deseja apagar permanentemente?").
				Value(&confirmDelete),
		),
	)

	if err := form.Run(); err != nil || !confirmDelete || selectedID == "" {
		tui.Info("Operação cancelada.")
		return
	}

	if err := config.DeleteClient(selectedID); err != nil {
		tui.Error(fmt.Sprintf("Erro ao remover: %v", err))
		return
	}

	tui.Success(fmt.Sprintf("Cliente '%s' removido com sucesso!", selectedID))
}
