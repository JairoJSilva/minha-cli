package cmd

import (
	"fmt"
	"os"

	"github.com/JairoJSilva/minha-cli/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mc",
	Short: "Minha-CLI - Multi-Cloud & SRE Context Manager",
	Long: `☁️  Minha-CLI (Go Edition)
Gerenciador de contextos de nuvem de alta performance para SREs e DevOps.
Permite alternar simultaneamente AWS, Oracle OCI, GCP, Azure e Kubernetes.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Se nenhum subcomando for passado, abre a TUI interativa
		selected, err := tui.RunMenu()
		if err != nil {
			fmt.Printf("Erro ao executar menu: %v\n", err)
			return
		}

		switch selected {
		case "switch":
			runSwitchInteractive()
		case "status":
			runStatus()
		case "test":
			runTestParallel()
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
		default:
			// Sair
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
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(k8sCmd)
}
