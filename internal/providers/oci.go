package providers

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type OCINamespaceResponse struct {
	Data string `json:"data"`
}

func TestOCI(profile string) (namespace, region string, err error) {
	if profile == "" {
		profile = "DEFAULT"
	}

	cmd := exec.Command("oci", "os", "ns", "get", "--profile", profile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("falha ao consultar Oracle OCI: %s", string(out))
	}

	var nsResp OCINamespaceResponse
	if err := json.Unmarshal(out, &nsResp); err == nil && nsResp.Data != "" {
		namespace = nsResp.Data
	} else {
		// Fallback de extração em texto
		namespace = profile
	}

	region = "sa-saopaulo-1"
	return namespace, region, nil
}
