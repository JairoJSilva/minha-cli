package cmd

import (
	"fmt"
	"os"

	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "mc",
	Version: AppVersion,
	Short:   "Minha-CLI - Multi-Cloud & SRE Context Manager",
	Long: `☁️  Minha-CLI (Go Edition)
Gerenciador de contextos de nuvem de alta performance para SREs e DevOps.
Permite alternar simultaneamente AWS, Oracle OCI, GCP, Azure e Kubernetes.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Banner exibido uma única vez na abertura
		tui.PrintBanner()

		// Loop principal: o menu reabre após cada ação
		for {
			selected, err := tui.RunMenuLoop()
			if err != nil {
				fmt.Printf("Erro ao executar menu: %v\n", err)
				return
			}

			// Sair: seleção vazia (q/ESC) ou opção "exit"
			if selected == "" || selected == "exit" {
				return
			}

			switch selected {
			case "switch":
				runSwitchInteractive()
			case "status":
				runStatus()
			case "test":
				runTestParallel()
			case "show":
				runShowInteractive()
			case "scan":
				runScan()
			case "add":
				runAddInteractive()
			case "edit":
				runEditInteractive()
			case "delete":
				runDeleteInteractive()
			case "list":
				runList()
			case "k8s":
				runK8s()
			case "clear":
				runClear()
			case "version":
				runVersion()
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(k8sCmd)
	rootCmd.AddCommand(versionCmd)
}
