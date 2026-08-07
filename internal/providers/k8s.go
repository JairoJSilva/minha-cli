package providers

import (
	"fmt"
	"os/exec"
	"strings"
)

func TestK8s(contextName string) (currentCtx string, nodeCount int, err error) {
	if contextName != "" {
		_ = SwitchK8sContext(contextName)
	}

	ctxOut, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		return "", 0, fmt.Errorf("nenhum contexto ativo")
	}
	currentCtx = strings.TrimSpace(string(ctxOut))

	// Conta nodes ativos
	nodeOut, err := exec.Command("kubectl", "get", "nodes", "--no-headers").Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(nodeOut)), "\n")
		for _, l := range lines {
			if strings.TrimSpace(l) != "" {
				nodeCount++
			}
		}
	}

	return currentCtx, nodeCount, nil
}

func SwitchK8sContext(contextName string) error {
	cmd := exec.Command("kubectl", "config", "use-context", contextName)
	return cmd.Run()
}
