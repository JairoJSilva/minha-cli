package env

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/JairoJSilva/minha-cli/internal/config"
)

type ActiveState struct {
	AWSProfile string
	OCIProfile string
	GCPConfig  string
	AzureSub   string
	K8sContext string
}

func GetCurrentState() ActiveState {
	aws := os.Getenv("AWS_PROFILE")
	if aws == "" {
		aws = "<não definido>"
	}

	oci := os.Getenv("OCI_CLI_PROFILE")
	if oci == "" {
		oci = "<não definido>"
	}

	gcp := os.Getenv("CLOUDSDK_ACTIVE_CONFIG_NAME")
	if gcp == "" {
		gcp = "<não definido>"
	}

	azure := os.Getenv("AZURE_SUBSCRIPTION")
	if azure == "" {
		azure = "<padrão/sessão>"
	}

	k8s := "<não conectado>"
	if out, err := exec.Command("kubectl", "config", "current-context").Output(); err == nil {
		ctx := strings.TrimSpace(string(out))
		if ctx != "" {
			k8s = ctx
		}
	}

	return ActiveState{
		AWSProfile: aws,
		OCIProfile: oci,
		GCPConfig:  gcp,
		AzureSub:   azure,
		K8sContext: k8s,
	}
}

// Gera o script shell a ser avaliado pelo wrapper (source <(mc env ...))
func GenerateExportScript(c *config.Client) string {
	var sb strings.Builder

	// AWS
	if c.AWSProfile != nil && *c.AWSProfile != "" {
		sb.WriteString(fmt.Sprintf("export AWS_PROFILE=\"%s\"\n", *c.AWSProfile))
	} else {
		sb.WriteString("unset AWS_PROFILE\n")
	}

	// OCI
	if c.OCIProfile != nil && *c.OCIProfile != "" {
		sb.WriteString(fmt.Sprintf("export OCI_CLI_PROFILE=\"%s\"\n", *c.OCIProfile))
	} else {
		sb.WriteString("unset OCI_CLI_PROFILE\n")
	}

	// GCP
	if c.GCPConfig != nil && *c.GCPConfig != "" {
		sb.WriteString(fmt.Sprintf("export CLOUDSDK_ACTIVE_CONFIG_NAME=\"%s\"\n", *c.GCPConfig))
	} else {
		sb.WriteString("unset CLOUDSDK_ACTIVE_CONFIG_NAME\n")
	}

	// Azure
	if c.AzureSub != nil && *c.AzureSub != "" {
		sb.WriteString(fmt.Sprintf("export AZURE_SUBSCRIPTION=\"%s\"\n", *c.AzureSub))
		sb.WriteString(fmt.Sprintf("if command -v az >/dev/null 2>&1; then az account set --subscription \"%s\" >/dev/null 2>&1 || true; fi\n", *c.AzureSub))
	} else {
		sb.WriteString("unset AZURE_SUBSCRIPTION\n")
	}

	// Kubernetes
	if c.K8sContext != nil && *c.K8sContext != "" {
		sb.WriteString(fmt.Sprintf("if command -v kubectl >/dev/null 2>&1; then kubectl config use-context \"%s\" >/dev/null 2>&1 || true; fi\n", *c.K8sContext))
	}

	return sb.String()
}

// Gera o script shell de limpeza total
func GenerateClearScript() string {
	return `unset AWS_PROFILE
unset AWS_ACCESS_KEY_ID
unset AWS_SECRET_ACCESS_KEY
unset AWS_SESSION_TOKEN
unset AWS_DEFAULT_REGION
unset OCI_CLI_PROFILE
unset OCI_CLI_REGION
unset CLOUDSDK_ACTIVE_CONFIG_NAME
unset GOOGLE_APPLICATION_CREDENTIALS
unset AZURE_SUBSCRIPTION
unset AZURE_SUBSCRIPTION_NAME
unset KUBECONFIG
`
}
