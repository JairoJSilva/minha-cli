package config

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type ScanResult struct {
	AWSProfiles   []string
	OCIProfiles   []string
	GCPConfigs    []string
	K8sContexts   []string
	AzureSubs     []string
	ImportedCount int
	ExistingCount int
}

func ScanLocalEnvironment() (*ScanResult, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	result := &ScanResult{}

	// 1. AWS Scanner (~/.aws/credentials e ~/.aws/config)
	awsCreds := filepath.Join(home, ".aws", "credentials")
	if f, err := os.Open(awsCreds); err == nil {
		defer f.Close()
		re := regexp.MustCompile(`^\[(.*)\]$`)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			matches := re.FindStringSubmatch(strings.TrimSpace(scanner.Text()))
			if len(matches) > 1 {
				p := strings.TrimSpace(matches[1])
				if p != "" && !contains(result.AWSProfiles, p) {
					result.AWSProfiles = append(result.AWSProfiles, p)
				}
			}
		}
	}

	awsConfig := filepath.Join(home, ".aws", "config")
	if f, err := os.Open(awsConfig); err == nil {
		defer f.Close()
		re := regexp.MustCompile(`^\[(?:profile\s+)?(.*)\]$`)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			matches := re.FindStringSubmatch(strings.TrimSpace(scanner.Text()))
			if len(matches) > 1 {
				p := strings.TrimSpace(matches[1])
				if p != "" && !contains(result.AWSProfiles, p) {
					result.AWSProfiles = append(result.AWSProfiles, p)
				}
			}
		}
	}

	// 2. OCI Scanner (~/.oci/config)
	ociConfig := filepath.Join(home, ".oci", "config")
	if f, err := os.Open(ociConfig); err == nil {
		defer f.Close()
		re := regexp.MustCompile(`^\[(.*)\]$`)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			matches := re.FindStringSubmatch(strings.TrimSpace(scanner.Text()))
			if len(matches) > 1 {
				p := strings.TrimSpace(matches[1])
				if p != "" && !contains(result.OCIProfiles, p) {
					result.OCIProfiles = append(result.OCIProfiles, p)
				}
			}
		}
	}

	// 3. GCP Scanner (~/.config/gcloud/configurations/)
	gcpDir := filepath.Join(home, ".config", "gcloud", "configurations")
	if entries, err := os.ReadDir(gcpDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasPrefix(e.Name(), "config_") {
				cname := strings.TrimPrefix(e.Name(), "config_")
				if cname != "" && !contains(result.GCPConfigs, cname) {
					result.GCPConfigs = append(result.GCPConfigs, cname)
				}
			}
		}
	}

	// 4. Kubernetes Scanner via kubectl
	if out, err := exec.Command("kubectl", "config", "get-contexts", "-o", "name").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, l := range lines {
			l = strings.TrimSpace(l)
			if l != "" && !contains(result.K8sContexts, l) {
				result.K8sContexts = append(result.K8sContexts, l)
			}
		}
	}

	// 5. Azure Scanner (~/.azure/azureProfile.json ou az CLI)
	if out, err := exec.Command("az", "account", "list", "--query", "[].name", "-o", "tsv").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, l := range lines {
			l = strings.TrimSpace(l)
			if l != "" && !contains(result.AzureSubs, l) {
				result.AzureSubs = append(result.AzureSubs, l)
			}
		}
	}

	// 6. Fusão inteligente no clients.json
	clients, _ := LoadClients()
	existingKeys := make(map[string]bool)
	for _, c := range clients {
		existingKeys[strings.ToLower(c.ID)] = true
		existingKeys[strings.ToLower(c.Name)] = true
	}

	// Lista de chaves descobertas
	var discoveredSlugs []string
	addSlug := func(s string) {
		slug := strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(s, ""))
		if slug != "" && slug != "default" && !contains(discoveredSlugs, slug) {
			discoveredSlugs = append(discoveredSlugs, slug)
		}
	}

	for _, p := range result.AWSProfiles {
		addSlug(p)
	}
	for _, p := range result.OCIProfiles {
		addSlug(p)
	}
	for _, p := range result.GCPConfigs {
		addSlug(p)
	}

	for _, slug := range discoveredSlugs {
		if existingKeys[slug] {
			result.ExistingCount++
			continue
		}

		// Procura correspondência
		var awsP, ociP, gcpP, k8sP *string

		for _, p := range result.AWSProfiles {
			if strings.Contains(strings.ToLower(p), slug) {
				awsP = StringPtr(p)
				break
			}
		}
		for _, p := range result.OCIProfiles {
			if strings.Contains(strings.ToLower(p), slug) {
				ociP = StringPtr(p)
				break
			}
		}
		for _, p := range result.GCPConfigs {
			if strings.Contains(strings.ToLower(p), slug) {
				gcpP = StringPtr(p)
				break
			}
		}
		for _, p := range result.K8sContexts {
			if strings.Contains(strings.ToLower(p), slug) {
				k8sP = StringPtr(p)
				break
			}
		}

		displayName := strings.Title(slug) + " (Importado)"
		newClient := Client{
			ID:         slug,
			Name:       displayName,
			AWSProfile: awsP,
			OCIProfile: ociP,
			GCPConfig:  gcpP,
			K8sContext: k8sP,
		}

		_ = AddClient(newClient)
		result.ImportedCount++
	}

	return result, nil
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}
