package cmd

import (
	"fmt"
	"os"
	"path/filepath"
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

func saveAWSCredential(profile, accessKey, secretKey string) error {
	if profile == "" || accessKey == "" || secretKey == "" {
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	awsDir := filepath.Join(home, ".aws")
	_ = os.MkdirAll(awsDir, 0700)
	credPath := filepath.Join(awsDir, "credentials")

	content := ""
	if data, err := os.ReadFile(credPath); err == nil {
		content = string(data)
	}

	// Se o profile já existe no credentials, avisa
	if strings.Contains(content, "["+profile+"]") {
		return nil
	}

	block := fmt.Sprintf("\n[%s]\naws_access_key_id = %s\naws_secret_access_key = %s\n", profile, accessKey, secretKey)
	f, err := os.OpenFile(credPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(block)
	return err
}

func runAddInteractive() {
	tui.PrintBanner()
	fmt.Println("\n➕ Cadastrar Nova Conta / Cliente")

	var name, id, aws, oci, gcp, azure, k8s string
	var configAWSKeys bool
	var awsAccessKey, awsSecretKey string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Nome de Exibição").
				Description("Ex: Hospital Albert Einstein, Santander, Dentalis").
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
				Placeholder("ex: einstein").
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

	// Se o cliente já existe, oferece edição imediata em vez de erro
	if existing, err := config.FindClient(id); err == nil && existing != nil {
		var wantEdit bool
		formRedirect := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("O cliente '%s' (%s) já está cadastrado. Deseja editá-lo agora?", existing.Name, id)).
					Value(&wantEdit),
			),
		)
		if err := formRedirect.Run(); err == nil && wantEdit {
			runEditClientByID(id)
			return
		}
		tui.Info("Operação finalizada.")
		return
	}

	// Se preencheu AWS profile, pergunta se deseja salvar chaves no ~/.aws/credentials
	if aws != "" && !checkAWSProfileExists(aws) {
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
				tui.Success(fmt.Sprintf("Chaves da AWS salvas com sucesso no ~/.aws/credentials para o perfil '%s'!", aws))
			}
		}
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
