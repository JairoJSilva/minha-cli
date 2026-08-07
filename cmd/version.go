package cmd

import (
	"fmt"
	"runtime"

	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/spf13/cobra"
)

const (
	AppVersion = "v2.0.1"
	AppBuild   = "Go Edition (SRE Native)"
	AppCommit  = "dab536d"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v", "-v", "--version"},
	Short:   "Exibe a versão e detalhes de compilação da CLI",
	Run: func(cmd *cobra.Command, args []string) {
		runVersion()
	},
}

func runVersion() {
	body := fmt.Sprintf(` Versão       : %s
 Build        : %s
 Commit       : %s
 Go Runtime   : %s (%s/%s)
 SRE Engine   : Multi-Cloud Context Orchestrator
 Repositório  : github.com/JairoJSilva/minha-cli`,
		AppVersion, AppBuild, AppCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH)

	tui.PrintBanner()
	tui.PrintCard("MINHA-CLI — INFORMAÇÕES DE VERSÃO", body)
}
