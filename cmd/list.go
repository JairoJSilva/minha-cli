package cmd

import (
	"fmt"
	"strings"

	"github.com/JairoJSilva/minha-cli/internal/config"
	"github.com/JairoJSilva/minha-cli/internal/providers"
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

func runList() {
	clients, err := config.LoadClients()
	if err != nil || len(clients) == 0 {
		tui.Warn("Nenhum cliente cadastrado. Use 'mc add' ou 'mc scan'.")
		return
	}

	var sb strings.Builder
	for _, c := range clients {
		aws := config.SafeString(c.AWSProfile)
		if aws == "" {
			aws = "-"
		}
		oci := config.SafeString(c.OCIProfile)
		if oci == "" {
			oci = "-"
		}
		gcp := config.SafeString(c.GCPConfig)
		if gcp == "" {
			gcp = "-"
		}
		k8s := config.SafeString(c.K8sContext)
		if k8s == "" {
			k8s = "-"
		}

		sb.WriteString(fmt.Sprintf(" • %s [ID: %s]\n", c.Name, c.ID))
		sb.WriteString(fmt.Sprintf("   └─ AWS: %s | OCI: %s | GCP: %s | K8s: %s\n", aws, oci, gcp, k8s))
	}

	tui.PrintCard("PERFIS CADASTRADOS (GO ENGINE)", strings.TrimRight(sb.String(), "\n"))
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
