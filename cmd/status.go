package cmd

import (
	"fmt"

	"github.com/JairoJSilva/minha-cli/internal/env"
	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st", "info"},
	Short:   "Exibe o status e as variáveis ativas no momento",
	Run: func(cmd *cobra.Command, args []string) {
		runStatus()
	},
}

func runStatus() {
	st := env.GetCurrentState()

	body := fmt.Sprintf(` AWS Profile  : %s
 OCI Profile  : %s
 GCP Config   : %s
 Azure Context: %s
 Kubernetes   : %s`,
		st.AWSProfile, st.OCIProfile, st.GCPConfig, st.AzureSub, st.K8sContext)

	tui.PrintCard("STATUS DO CONTEXTO ATIVO", body)
}
