package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/JairoJSilva/minha-cli/internal/config"
	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/JairoJSilva/minha-cli/internal/vault"
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

	// ── Etapa 1: Metadados ──────────────────────────────────────────────────
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

	// Se o cliente já existe, oferece edição imediata
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

	// ── Etapa 2: Credenciais no Vault ────────────────────────────────────────
	var saveToVault bool
	formVault := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("🔐 Deseja salvar as credenciais de acesso neste cliente no Vault seguro?").
				Description("As chaves ficam criptografadas em ~/.config/minha-cli/vault.enc").
				Value(&saveToVault),
		),
	)
	_ = formVault.Run()

	hasVaultSecret := false

	if saveToVault {
		secret := vault.VaultSecret{}

		// AWS Credentials
		if aws != "" {
			var awsKey, awsSecret, awsRegion string
			var awsAuthMode string

			formAWSMode := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Modo de autenticação AWS").
						Options(
							huh.NewOption("🔑 Profile (Access Key + Secret)", "profile"),
							huh.NewOption("🎫 STS Token Temporário (AssumeRole/MFA)", "sts"),
						).
						Value(&awsAuthMode),
				),
			)
			if err := formAWSMode.Run(); err == nil {
				formAWSCreds := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("AWS Access Key ID").
							Placeholder("AKIA...").
							Value(&awsKey),
						huh.NewInput().
							Title("AWS Secret Access Key").
							EchoMode(huh.EchoModePassword).
							Value(&awsSecret),
						huh.NewInput().
							Title("AWS Region").
							Placeholder("us-east-1").
							Value(&awsRegion),
					),
				)
				if err := formAWSCreds.Run(); err == nil && awsKey != "" {
					secret.AWSAccessKeyID = awsKey
					secret.AWSSecretAccessKey = awsSecret
					secret.AWSRegion = awsRegion
					if awsAuthMode == "sts" {
						var roleARN string
						formRole := huh.NewForm(
							huh.NewGroup(
								huh.NewInput().
									Title("ARN do Role para AssumeRole (opcional)").
									Placeholder("arn:aws:iam::123456789012:role/MyRole").
									Value(&roleARN),
							),
						)
						_ = formRole.Run()
						secret.AWSRoleARN = roleARN
					}
				}
			}
		}

		// OCI Credentials
		if oci != "" {
			var ociUser, ociTenancy, ociFingerprint, ociKeyPath, ociRegion string
			formOCI := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("OCI User OCID").Placeholder("ocid1.user.oc1..").Value(&ociUser),
					huh.NewInput().Title("OCI Tenancy OCID").Placeholder("ocid1.tenancy.oc1..").Value(&ociTenancy),
					huh.NewInput().Title("OCI Fingerprint").Placeholder("xx:xx:xx:...").Value(&ociFingerprint),
					huh.NewInput().Title("Caminho da Chave Privada OCI").Placeholder("~/.oci/oci_api_key.pem").Value(&ociKeyPath),
					huh.NewInput().Title("OCI Region").Placeholder("sa-saopaulo-1").Value(&ociRegion),
				),
			)
			if err := formOCI.Run(); err == nil && ociUser != "" {
				secret.OCIUserOCID = ociUser
				secret.OCITenancyOCID = ociTenancy
				secret.OCIFingerprint = ociFingerprint
				secret.OCIPrivateKeyPath = ociKeyPath
				secret.OCIRegion = ociRegion
			}
		}

		// Salva no vault
		if err := vault.Store(id, secret); err != nil {
			tui.Error(fmt.Sprintf("Falha ao salvar no vault: %v", err))
		} else {
			hasVaultSecret = true
			tui.Success(fmt.Sprintf("Credenciais salvas com segurança no vault! (%s)", vault.VaultFilePath()))
		}
	}

	// ── Etapa 3: Salva metadados no clients.json ─────────────────────────────
	client := config.Client{
		ID:             id,
		Name:           name,
		AWSProfile:     config.StringPtr(aws),
		OCIProfile:     config.StringPtr(oci),
		GCPConfig:      config.StringPtr(gcp),
		AzureSub:       config.StringPtr(azure),
		K8sContext:     config.StringPtr(k8s),
		HasVaultSecret: hasVaultSecret,
	}

	if err := config.AddClient(client); err != nil {
		tui.Error(fmt.Sprintf("Falha ao salvar cliente: %v", err))
		return
	}

	tui.Success(fmt.Sprintf("Cliente '%s' (%s) adicionado com sucesso!", name, id))
}
