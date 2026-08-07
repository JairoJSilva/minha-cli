package providers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type AWSCallerIdentity struct {
	Account string `json:"Account"`
	Arn     string `json:"Arn"`
	UserId  string `json:"UserId"`
}

func TestAWS(profile string) (*AWSCallerIdentity, string, error) {
	if profile == "" {
		profile = "default"
	}

	var cmd *exec.Cmd
	if _, err := exec.LookPath("rtk"); err == nil {
		cmd = exec.Command("rtk", "aws", "sts", "get-caller-identity", "--profile", profile, "--output", "json")
	} else if _, err := exec.LookPath("/root/.local/bin/rtk"); err == nil {
		cmd = exec.Command("/root/.local/bin/rtk", "aws", "sts", "get-caller-identity", "--profile", profile, "--output", "json")
	} else {
		cmd = exec.Command("aws", "sts", "get-caller-identity", "--profile", profile, "--output", "json")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", fmt.Errorf("falha ao consultar AWS STS: %s", string(out))
	}

	var identity AWSCallerIdentity
	if err := json.Unmarshal(out, &identity); err != nil {
		// Tenta parser de fallback se for output direto
		if strings.Contains(string(out), "Account") {
			return &AWSCallerIdentity{Account: "Validado", Arn: string(out)}, "us-east-1", nil
		}
		return nil, "", fmt.Errorf("resposta inesperada da AWS: %s", string(out))
	}

	// Região configurada
	regCmd := exec.Command("aws", "configure", "get", "region", "--profile", profile)
	regOut, _ := regCmd.Output()
	region := strings.TrimSpace(string(regOut))
	if region == "" {
		region = "us-east-1"
	}

	return &identity, region, nil
}
