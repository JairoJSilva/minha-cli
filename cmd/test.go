package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/JairoJSilva/minha-cli/internal/providers"
	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"whoami", "t", "validar"},
	Short:   "Testa e valida as credenciais ativas nas APIs em paralelo com Goroutines",
	Run: func(cmd *cobra.Command, args []string) {
		runTestParallel()
	},
}

func runTestParallel() {
	tui.PrintBanner()
	fmt.Println("\n🔍 Testando Conexões e Identidades em Paralelo (Goroutines)...")

	var wg sync.WaitGroup

	// AWS
	wg.Add(1)
	go func() {
		defer wg.Done()
		currentAWS := os.Getenv("AWS_PROFILE")
		identity, region, err := providers.TestAWS(currentAWS)
		if err == nil && identity != nil {
			body := fmt.Sprintf(` Conta ID : %s
 Região   : %s
 Identity : %s
 Engine   : Go Nativo (Goroutines)
 Status   : AUTENTICADO & OPERACIONAL`, identity.Account, region, identity.Arn)
			tui.PrintCard("☁️  AWS - IDENTIDADE VALIDADA", body)
		} else {
			tui.Warn(fmt.Sprintf("AWS STS: não autenticado no perfil '%s'", currentAWS))
		}
	}()

	// Oracle OCI
	wg.Add(1)
	go func() {
		defer wg.Done()
		currentOCI := os.Getenv("OCI_CLI_PROFILE")
		ns, region, err := providers.TestOCI(currentOCI)
		if err == nil && ns != "" {
			body := fmt.Sprintf(` Namespace : %s
 Profile   : %s
 Região    : %s
 Status    : AUTENTICADO & OPERACIONAL`, ns, currentOCI, region)
			tui.PrintCard("🏛️  ORACLE CLOUD (OCI) - IDENTIDADE VALIDADA", body)
		} else {
			tui.Warn(fmt.Sprintf("OCI: perfil '%s' não respondeu ou não configurado", currentOCI))
		}
	}()

	// Kubernetes
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, nodeCount, err := providers.TestK8s("")
		if err == nil && ctx != "" {
			body := fmt.Sprintf(` Contexto    : %s
 Nodes Ativos: %d
 Status      : CLUSTER ONLINE & RESPONSIVO`, ctx, nodeCount)
			tui.PrintCard("☸️  KUBERNETES - CLUSTER ATIVO", body)
		}
	}()

	wg.Wait()
	fmt.Println()
}
