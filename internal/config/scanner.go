package config

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type AWSProfileDetail struct {
	Name      string
	AccessKey string
	Region    string
	Source    string
}

type OCIProfileDetail struct {
	Name    string
	Tenancy string
	Region  string
}

type K8sContextDetail struct {
	Name      string
	Cluster   string
	Namespace string
}

type GCPConfigDetail struct {
	Name    string
	Project string
	Account string
}

type AzureSubDetail struct {
	Name string
	ID   string
}

type ScanReport struct {
	AWSDetails      []AWSProfileDetail
	OCIDetails      []OCIProfileDetail
	GCPDetails      []GCPConfigDetail
	K8sDetails      []K8sContextDetail
	AzureDetails    []AzureSubDetail
	ImportedCount   int
	ExistingCount   int
	TotalRegistered int
	RegisteredNames []string
}

func maskSecret(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "..." + s[len(s)-4:]
}

func ScanLocalEnvironmentDetailed() (*ScanReport, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	report := &ScanReport{}

	// 1. Scanner AWS Detalhado (~/.aws/credentials e ~/.aws/config)
	awsCreds := filepath.Join(home, ".aws", "credentials")
	if f, err := os.Open(awsCreds); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		var currentProfile *AWSProfileDetail

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				if currentProfile != nil {
					report.AWSDetails = append(report.AWSDetails, *currentProfile)
				}
				pName := strings.Trim(line, "[]")
				currentProfile = &AWSProfileDetail{Name: pName, Source: "~/.aws/credentials", Region: "us-east-1"}
			} else if currentProfile != nil {
				if strings.HasPrefix(line, "aws_access_key_id") {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						currentProfile.AccessKey = maskSecret(strings.TrimSpace(parts[1]))
					}
				}
			}
		}
		if currentProfile != nil {
			report.AWSDetails = append(report.AWSDetails, *currentProfile)
		}
	}

	// 2. Scanner Oracle OCI Detalhado (~/.oci/config)
	ociConfig := filepath.Join(home, ".oci", "config")
	if f, err := os.Open(ociConfig); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		var currentOCI *OCIProfileDetail

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				if currentOCI != nil {
					report.OCIDetails = append(report.OCIDetails, *currentOCI)
				}
				pName := strings.Trim(line, "[]")
				currentOCI = &OCIProfileDetail{Name: pName, Region: "sa-saopaulo-1"}
			} else if currentOCI != nil {
				if strings.HasPrefix(line, "tenancy=") {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						currentOCI.Tenancy = maskSecret(strings.TrimSpace(parts[1]))
					}
				} else if strings.HasPrefix(line, "region=") {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						currentOCI.Region = strings.TrimSpace(parts[1])
					}
				}
			}
		}
		if currentOCI != nil {
			report.OCIDetails = append(report.OCIDetails, *currentOCI)
		}
	}

	// 3. Scanner Kubernetes Detalhado (kubectl)
	if out, err := exec.Command("kubectl", "config", "get-contexts", "--no-headers").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, l := range lines {
			fields := strings.Fields(l)
			if len(fields) >= 2 {
				name := fields[0]
				if name == "*" && len(fields) >= 3 {
					name = fields[1]
				}
				cluster := fields[len(fields)-2]
				report.K8sDetails = append(report.K8sDetails, K8sContextDetail{
					Name:    name,
					Cluster: cluster,
				})
			}
		}
	}

	// 4. Scanner Google Cloud (~/.config/gcloud/configurations/)
	gcpDir := filepath.Join(home, ".config", "gcloud", "configurations")
	if entries, err := os.ReadDir(gcpDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasPrefix(e.Name(), "config_") {
				cname := strings.TrimPrefix(e.Name(), "config_")
				if cname != "" {
					report.GCPDetails = append(report.GCPDetails, GCPConfigDetail{Name: cname})
				}
			}
		}
	}

	// 5. Scanner Azure
	if out, err := exec.Command("az", "account", "list", "--query", "[].{name:name, id:id}", "-o", "tsv").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, l := range lines {
			parts := strings.Split(l, "\t")
			if len(parts) >= 2 {
				report.AzureDetails = append(report.AzureDetails, AzureSubDetail{
					Name: parts[0],
					ID:   maskSecret(parts[1]),
				})
			}
		}
	}

	// 6. Carrega os clientes já cadastrados no clients.json
	clients, _ := LoadClients()
	report.TotalRegistered = len(clients)
	existingKeys := make(map[string]bool)
	for _, c := range clients {
		existingKeys[strings.ToLower(c.ID)] = true
		existingKeys[strings.ToLower(c.Name)] = true
		report.RegisteredNames = append(report.RegisteredNames, fmt.Sprintf("%s [%s]", c.Name, c.ID))
	}
	report.ExistingCount = len(clients)

	// 7. Auto-importação de novos perfis encontrados
	for _, awsP := range report.AWSDetails {
		slug := strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(awsP.Name, ""))
		if slug == "" || slug == "default" {
			continue
		}

		if !existingKeys[slug] {
			displayName := strings.Title(slug) + " (Importado)"
			_ = AddClient(Client{
				ID:         slug,
				Name:       displayName,
				AWSProfile: StringPtr(awsP.Name),
			})
			report.ImportedCount++
			report.RegisteredNames = append(report.RegisteredNames, fmt.Sprintf("%s [%s]", displayName, slug))
			report.TotalRegistered++
			existingKeys[slug] = true
		}
	}

	return report, nil
}
