package cmd

import (
	"fmt"

	"github.com/JairoJSilva/minha-cli/internal/config"
	"github.com/JairoJSilva/minha-cli/internal/env"
	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/JairoJSilva/minha-cli/internal/vault"
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

	// ── Modo Vault: credenciais armazenadas de forma segura ──────────────────
	if client.HasVaultSecret {
		secret, err := vault.Get(client.ID)
		if err != nil {
			tui.Error(fmt.Sprintf("Falha ao ler credenciais do vault: %v", err))
			tui.Warn("Tentando fallback para variáveis de ambiente...")
		} else {
			// Gera script com as credenciais reais do vault (nunca escritas em disco)
			script := env.GenerateExportScriptFromVault(client, secret)
			env.WriteEnvToFile(script)

			// Feedback
			tui.Success(fmt.Sprintf("🔐 Contexto ativado via Vault: %s", client.Name))
			aws := "<nenhum>"
			if secret.AWSAccessKeyID != "" {
				aws = "🔐 via vault"
			} else if client.AWSProfile != nil {
				aws = *client.AWSProfile
			}
			oci := "<nenhum>"
			if secret.OCIUserOCID != "" {
				oci = "🔐 via vault"
			} else if client.OCIProfile != nil {
				oci = *client.OCIProfile
			}
			gcp := config.SafeString(client.GCPConfig)
			if gcp == "" {
				gcp = "<nenhum>"
			}
			k8s := config.SafeString(client.K8sContext)
			if k8s == "" {
				k8s = "<nenhum>"
			}
			fmt.Printf("  \033[2mAWS: %s | OCI: %s | GCP: %s | K8s: %s\033[0m\n", aws, oci, gcp, k8s)
			fmt.Printf("  \033[2m⏱️  Sessão ativa — variáveis serão removidas ao fechar o terminal\033[0m\n")
			return
		}
	}

	// ── Modo Legacy: sem vault (profile-based) ───────────────────────────────
	script := env.GenerateExportScript(client)
	env.WriteEnvToFile(script)

	tui.Success(fmt.Sprintf("Contexto ativado: %s", client.Name))
	aws := "<nenhum>"
	if client.AWSProfile != nil {
		aws = *client.AWSProfile
		if !config.AWSProfileLocalExists(aws) {
			fmt.Printf("  \033[33m⚠️  Aviso: O profile AWS '%s' não está no ~/.aws/credentials.\033[0m\n", aws)
			fmt.Printf("  \033[33m   → Use 'mc add' para salvar as chaves no vault seguro.\033[0m\n")
		}
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
		label := c.Name
		if c.HasVaultSecret {
			label = "🔐 " + label
		}
		options = append(options, huh.NewOption(label, c.ID))
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
