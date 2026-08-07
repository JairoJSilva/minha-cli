package providers

import (
	"os/exec"
	"strings"
)

func TestGCP() (project, account string) {
	projOut, _ := exec.Command("gcloud", "config", "get-value", "project").Output()
	project = strings.TrimSpace(string(projOut))

	accOut, _ := exec.Command("gcloud", "config", "get-value", "account").Output()
	account = strings.TrimSpace(string(accOut))

	return project, account
}
