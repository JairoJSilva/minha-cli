package cmd

import (
	"fmt"
	"strings"

	"github.com/JairoJSilva/minha-cli/internal/config"
	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "listar"},
	Short:   "Lista todas as contas e nuvens cadastradas",
	Run: func(cmd *cobra.Command, args []string) {
		runList()
	},
}

// credStatus retorna o indicador visual para uma credencial
func credStatus(profileName string, hasVault bool, localCheckFn func(string) bool) string {
	if profileName == "" {
		return "-"
	}
	if hasVault {
		return "🔐 " + profileName
	}
	if localCheckFn(profileName) {
		return "✅ " + profileName
	}
	return "📋 " + profileName + " (sem credencial)"
}

func runList() {
	clients, err := config.LoadClients()
	if err != nil || len(clients) == 0 {
		tui.Warn("Nenhum cliente cadastrado. Use 'mc add' ou 'mc scan'.")
		return
	}

	var sb strings.Builder
	for _, c := range clients {
		vaultTag := ""
		if c.HasVaultSecret {
			vaultTag = " 🔐"
		}

		sb.WriteString(fmt.Sprintf(" • %s [ID: %s]%s\n", c.Name, c.ID, vaultTag))

		// AWS
		if aws := config.SafeString(c.AWSProfile); aws != "" {
			sb.WriteString(fmt.Sprintf("   ├─ AWS  : %s\n",
				credStatus(aws, c.HasVaultSecret, config.AWSProfileLocalExists)))
		}

		// OCI
		if oci := config.SafeString(c.OCIProfile); oci != "" {
			sb.WriteString(fmt.Sprintf("   ├─ OCI  : %s\n",
				credStatus(oci, c.HasVaultSecret, config.OCIProfileLocalExists)))
		}

		// GCP
		if gcp := config.SafeString(c.GCPConfig); gcp != "" {
			sb.WriteString(fmt.Sprintf("   ├─ GCP  : %s\n",
				credStatus(gcp, c.HasVaultSecret, config.GCPConfigLocalExists)))
		}

		// Azure (só referência, sem checagem local possível)
		if azure := config.SafeString(c.AzureSub); azure != "" {
			sb.WriteString(fmt.Sprintf("   ├─ Azure: %s\n", azure))
		}

		// K8s
		if k8s := config.SafeString(c.K8sContext); k8s != "" {
			sb.WriteString(fmt.Sprintf("   └─ K8s  : %s\n",
				credStatus(k8s, false, config.K8sContextLocalExists)))
		}
	}

	legend := "\n  ✅ credencial local OK  📋 só referência (sem credencial)  🔐 no vault"
	tui.PrintCard("PERFIS CADASTRADOS", strings.TrimRight(sb.String(), "\n")+"\n"+legend)
}

var k8sCmd = &cobra.Command{
	Use:     "k8s",
	Aliases: []string{"kube", "cluster"},
	Short:   "Verifica o status e nós do cluster Kubernetes ativo",
	Run: func(cmd *cobra.Command, args []string) {
		runK8s()
	},
}

func runK8s() {
	ctx, nodeCount, err := providers.TestK8s("")
	if err != nil || ctx == "" {
		tui.Warn("Não foi possível conectar ao Kubernetes ativo.")
		return
	}

	body := fmt.Sprintf(` Contexto Ativo : %s
 Nodes Ativos   : %d
 Status         : CLUSTER ONLINE & RESPONSIVO`, ctx, nodeCount)

	tui.PrintCard("☸️  KUBERNETES - CLUSTER ATIVO", body)
}
