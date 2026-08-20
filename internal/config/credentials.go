package config

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// AWSProfileLocalExists verifica se um profile AWS existe em ~/.aws/credentials ou ~/.aws/config
func AWSProfileLocalExists(profile string) bool {
	if profile == "" || profile == "default" {
		return true
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	for _, file := range []string{
		filepath.Join(home, ".aws", "credentials"),
		filepath.Join(home, ".aws", "config"),
	} {
		if data, err := os.ReadFile(file); err == nil {
			content := string(data)
			if strings.Contains(content, "["+profile+"]") ||
				strings.Contains(content, "[profile "+profile+"]") {
				return true
			}
		}
	}
	return false
}

// OCIProfileLocalExists verifica se um profile OCI existe em ~/.oci/config
func OCIProfileLocalExists(profile string) bool {
	if profile == "" {
		return false
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	ociConfig := filepath.Join(home, ".oci", "config")
	f, err := os.Open(ociConfig)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "["+profile+"]" || line == "["+strings.ToUpper(profile)+"]" {
			return true
		}
	}
	return false
}

// GCPConfigLocalExists verifica se uma configuração GCP existe em ~/.config/gcloud/configurations/
func GCPConfigLocalExists(configName string) bool {
	if configName == "" {
		return false
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	configFile := filepath.Join(home, ".config", "gcloud", "configurations", "config_"+configName)
	_, err = os.Stat(configFile)
	return err == nil
}

// K8sContextLocalExists verifica se um contexto Kubernetes existe no kubeconfig
func K8sContextLocalExists(ctxName string) bool {
	if ctxName == "" {
		return false
	}

	out, err := exec.Command("kubectl", "config", "get-contexts", "-o", "name").Output()
	if err != nil {
		return false
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.TrimSpace(line) == ctxName {
			return true
		}
	}
	return false
}
